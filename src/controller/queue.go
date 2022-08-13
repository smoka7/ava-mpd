package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
)

func (c Mpd) Queue(w http.ResponseWriter, r *http.Request) {
	c.Client.Connect()
	defer c.Client.Close()

	switch r.Method {
	case http.MethodPost:
		var request playlist.QueueRequest
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		c.runCommand(request)

	case http.MethodGet:
		err = r.ParseForm()
		config.Log(err)

		playlist := playlist.GetQueue(c.Client, r.Form.Get("page"))
		err = json.NewEncoder(w).Encode(playlist)
		config.Log(err)
	}
}

type Commands map[string]func()

func (c Mpd) runCommand(request playlist.QueueRequest) {
	a := playlist.NewAction(c.Client)
	cmd := Commands{
		"play":         func() { a.PlaySong(request.Data) },
		"delete":       func() { a.DeleteSong(request.Data) },
		"move":         func() { a.MoveSong(request.Data) },
		"save":         func() { a.SaveQueue(request.Data) },
		"shuffleAlbum": func() { a.ShuffleAlbum(request.Data) },
		"addsong": func() {
			file, err := a.GetSongFile(request.Data.ID)
			if err != nil {
				return
			}
			a.AddSongToPlaylist(request.Data.Playlist, file)
		},
		"removeduplicate": func() {
			// TODO
			a.RemoveDuplicatesongs("")
		},
	}
	if action, ok := cmd[request.Command]; ok {
		action()
	}
}
