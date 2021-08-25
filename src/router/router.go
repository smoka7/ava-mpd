package router

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/controller"
	"github.com/smoka7/ava/src/watcher"
	"golang.org/x/net/websocket"
)

var cl controller.Mpd
var wc watcher.Mpd

func init() {
	cl.Client.ReadEnv()
	cl.Client.ReadFromFlags()
	wc.Client.Address = cl.Client.Address
	wc.Client.Password = cl.Client.Password
	ip := getHostIP()
	fmt.Printf("\n--- serving on %s:%s\n", ip, cl.Client.AppPort)
	serveCoverArts()
	go wc.RecordPlayCount()
}

//sets cache-control max-age's to 2 days
func CacheControl(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=172800")
		h.ServeHTTP(w, r)
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
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
