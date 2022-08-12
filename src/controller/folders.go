package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
)

type FolderRequest struct {
	Command string     `json:"command"`
	Data    FolderData `json:"data"`
}

type FolderData struct {
	Pos  string
	File string
}

func (c *Mpd) ServerFolders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var request FolderRequest
		c.Client.Connect()
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "add":
			playlist.AddFolder(&c.Client, request.Data.Pos, request.Data.File)
		case "addToPlaylist":
			// TODO
			// playlist.AddSongToPlaylist(&c.Client, request.Data.Pos, request.Data.File)
		case "play":
			playlist.PlayFolder(&c.Client, request.Data.File)
		case "list":
			folders := playlist.ListFolders(&c.Client, request.Data.File)
			err = json.NewEncoder(w).Encode(folders)
			config.Log(err)
		}
		c.Client.Close()
		return
	}
}
