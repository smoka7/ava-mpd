package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
)

func (c Mpd) StoredPlaylist(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a := playlist.NewAction(c.Client)
		playlist := a.ListStoredPlaylist()
		err = json.NewEncoder(w).Encode(playlist)
		config.Log(err)

	case http.MethodPost:
		var request playlist.PlaylistRequest
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)

		c.Client.Connect()
		defer c.Client.Close()
		c.runPlCommand(w, request)
	}
}

func (c Mpd) runPlCommand(w http.ResponseWriter, request playlist.PlaylistRequest) {
	a := playlist.NewAction(c.Client)
	cmd := Commands{
		"load":            func() { a.LoadPlaylist(request.Data) },
		"clear":           func() { a.ClearPlaylist(request.Data.Playlist) },
		"delete":          func() { a.DeletePlaylist(request.Data.Playlist) },
		"play":            func() { a.PlayPlaylist(request.Data) },
		"removeduplicate": func() { a.RemoveDuplicatesongs(request.Data.Playlist) },
		"removeinvalid":   func() { a.RemoveInvalidsongs(request.Data.Playlist) },
		"rename":          func() { a.RenamePlaylist(request.Data) },
		"playliked": func() {
			songs := a.GetLikedSongs()
			a.PlayFolder(songs...)
		},
		"addliked": func() {
			songs := a.GetLikedSongs()
			a.AddFolder(request.Data.Pos, songs...)
		},
		"playmost": func() {
			songs := a.GetMostPlayed()
			a.PlayFolder(songs...)
		},
		"addmost": func() {
			songs := a.GetMostPlayed()
			a.AddFolder(request.Data.Pos, songs...)
		},
		"list": func() {
			songs := a.ListSongsIn(request.Data.Playlist)
			err = json.NewEncoder(w).Encode(songs)
			config.Log(err)
		},
	}

	if action, ok := cmd[request.Command]; ok {
		action()
	}
}
