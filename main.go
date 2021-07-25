package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/router"
)

//go:embed front/dist/*
var front embed.FS

func getHostIP() net.IP {
	conn, err := net.Dial("udp", "192.168.1.1:80")
	config.Log(err)
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

//embeds frontEnd to the bin
func serveFrontEnd() {
	dist, err := fs.Sub(front, "front/dist")
	config.Log(err)
	frontEnd := http.FileServer(http.FS(dist))
	http.Handle("/", router.CacheControl(frontEnd))
}

func main() {
	ip := getHostIP()
	appPort := flag.String("c", "3001", "The port to run this app on it")
	flag.Parse()
	fmt.Printf("\n--- serving on %s:%s\n", ip, *appPort)
	serveFrontEnd()
	router.Router(*appPort)
	log.Panic(
		http.ListenAndServe(":"+*appPort, nil),
	)
}
