package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/playlist"
)

func (c *Mpd) Queue(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&request)
		switch request.Command {
		case "addsong":
			playlist.AddSongToPlaylist(c.Client, request.Data.Playlist, request.Data.Song)
		case "delete":
			playlist.DeleteSong(c.Client, request.Data.Start, -1)
		case "play":
			playlist.PlaySong(c.Client, request.Data.Start)
		case "save":
			playlist.SaveQueue(c.Client, request.Data.Playlist)
		}
		return
	}
	playlist := playlist.GetQueue(c.Client)
	json.NewEncoder(w).Encode(playlist)
}
