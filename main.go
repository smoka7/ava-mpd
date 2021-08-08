package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/router"
)

//go:embed front/dist/*
var front embed.FS

//embeds frontEnd to the bin
func serveFrontEnd() {
	dist, err := fs.Sub(front, "front/dist")
	config.Log(err)
	frontEnd := http.FileServer(http.FS(dist))
	http.Handle("/", router.CacheControl(frontEnd))
}

func main() {
	serveFrontEnd()
	router.Router()
}
