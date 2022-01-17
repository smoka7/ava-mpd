package router

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/controller"
	"github.com/smoka7/ava/src/watcher"
	"golang.org/x/net/websocket"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

var cl controller.Mpd
var wc watcher.Mpd

func init() {
	cl.Client.ReadEnv()
	cl.Client.ReadFromFlags()
	if cl.Client.Address == "" {
		log.Println("mpd server address is empty")
		os.Exit(1)
	}
	wc.Client.Address = cl.Client.Address
	wc.Client.Password = cl.Client.Password
	ip := getHostIP()
	fmt.Printf("\n--- serving on %s:%s\n", ip, cl.Client.AppPort)
	serveCoverArts()
	go wc.RecordPlayCount()
}

func Router() {
	http.Handle("/update", websocket.Handler(wc.Serve))
	http.HandleFunc("/api/playback", cl.Playback)
	http.HandleFunc("/api/status", cl.Status)
	http.HandleFunc("/api/stored", cl.StoredPlaylist)
	http.HandleFunc("/api/queue", cl.Queue)
	http.HandleFunc("/api/folders", cl.ServerFolders)
	http.HandleFunc("/api/search", cl.SearchServer)
	http.HandleFunc("/api/song", cl.Song)
	http.HandleFunc("/api/setting", cl.Settings)
	log.Panic(
		http.ListenAndServe(":"+cl.Client.AppPort, nil),
	)
}

// implements http.ResponseWriter
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//sets cache-control max-age's to 2 days and writes response in gzip if it is supported
func CacheControl(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=172800")
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

//returns caver art files from cache
func serveCoverArts() {
	cache, err := os.UserCacheDir()
	config.Log(err)
	coverArtDir := cache + "/ava-mpd"
	err = os.MkdirAll(coverArtDir, 0777)
	config.Log(err)
	http.Handle("/coverart/", CacheControl(http.FileServer(http.Dir(coverArtDir))))
}

//returns current host IP in LAN
func getHostIP() net.IP {
	conn, err := net.Dial("udp", "192.168.1.1:80")
	config.Log(err)
	if err != nil {
		return nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
