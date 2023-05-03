package router

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/controller"
	"github.com/smoka7/ava/src/watcher"
	"golang.org/x/net/websocket"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

var (
	cl controller.Mpd
	wc watcher.Mpd
)

func InitConnections(frontDist fs.FS, config config.Connection) {
	cl = controller.NewClient(config)
	wc = watcher.NewClient(config)

	ip := getHostIP()
	fmt.Printf("\n--- serving on %s:%s\n", ip, cl.AppPort)

	go wc.RecordPlayCount()

	frontEnd := http.FileServer(http.FS(frontDist))
	http.Handle("/", cacheControl(frontEnd))

	serveCoverArts()
	serveRoutes()

	srv := http.Server{
		Addr:         ":" + cl.AppPort,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Panic(srv.ListenAndServe())
}

func serveRoutes() {
	http.Handle("/update", websocket.Handler(wc.Serve))
	http.HandleFunc("/api/playback", cl.Playback)
	http.HandleFunc("/api/status", cl.Status)
	http.HandleFunc("/api/stored", cl.StoredPlaylist)
	http.HandleFunc("/api/queue", cl.Queue)
	http.HandleFunc("/api/folders", cl.ServerFolders)
	http.HandleFunc("/api/search", cl.SearchServer)
	http.HandleFunc("/api/song", cl.Song)
	http.HandleFunc("/api/setting", cl.Settings)
}

// implements http.ResponseWriter.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// sets cache-control max-age's to 2 days and writes response in gzip if it is supported.
func cacheControl(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=172800")

		// do not compress if client doesn't support it
		if !strings.Contains(r.Header.Get("accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)

			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		gzipResponse := gzipResponseWriter{Writer: gzipWriter, ResponseWriter: w}
		h.ServeHTTP(gzipResponse, r)
	}

	return http.HandlerFunc(fn)
}

// returns caver art files from cache.
func serveCoverArts() {
	cache, err := os.UserCacheDir()
	config.Log(err)
	coverArtDir := cache + "/ava-mpd"
	err = os.MkdirAll(coverArtDir, 0o777)
	config.Log(err)
	http.Handle("/coverart/", cacheControl(http.FileServer(http.Dir(coverArtDir))))
}

// returns current host IP in LAN.
func getHostIP() string {
	conn, err := net.Dial("udp", "192.168.1.1:80")
	config.Log(err)
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "127.0.0.1"
	}

	return localAddr.IP.String()
}
