package main

import (
	"embed"
	"io/fs"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/router"
)

//go:embed front/dist/*
var front embed.FS

func main() {
	frontDist, err := fs.Sub(front, "front/dist")
	config.Log(err)
	router.Router(&frontDist)
}
