package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
)

func (c Mpd) ServerFolders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var request playlist.FolderRequest
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)

		c.Connect()
		defer c.Close()
		c.runFolderCommand(w, request)
	}
}

func (c Mpd) runFolderCommand(w http.ResponseWriter, request playlist.FolderRequest) {
	a := playlist.NewAction(c.Connection)
	folderCmds := Commands{
		"add": func() {
			a.AddFolder(request.Data.Pos, request.Data.File)
		},
		"addToPlaylist": func() {
			a.AddSongToPlaylist(request.Data.Pos, request.Data.File)
		},
		"play": func() {
			a.PlayFolder(request.Data.File)
		},
		"list": func() {
			// TODO change to get
			folders := a.ListFolders(request.Data)
			err = json.NewEncoder(w).Encode(folders)
			config.Log(err)
		},
	}
	if action, ok := folderCmds[request.Command]; ok {
		action()
	}
}
