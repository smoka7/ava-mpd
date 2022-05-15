package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
)

func (c *Mpd) StoredPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.Client.Connect()
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "load":
			playlist.LoadPlaylist(&c.Client, request.Data.Playlist)
		case "clear":
			playlist.ClearPlaylist(&c.Client, request.Data.Playlist)
		case "delete":
			playlist.DeletePlaylist(&c.Client, request.Data.Playlist)
		case "play":
			playlist.PlayPlaylist(&c.Client, request.Data.Playlist)
		case "addafter":
			playlist.AddAfterCurrent(&c.Client, request.Data.Playlist)
		case "removeduplicate":
			playlist.RemoveDuplicatesongs(&c.Client, request.Data.Playlist)
		case "removeinvalid":
			playlist.RemoveInvalidsongs(&c.Client, request.Data.Playlist)
		case "rename":
			playlist.RenamePlaylist(&c.Client, request.Data.Playlist, request.Data.Song)
		case "playliked":
			songs := playlist.GetLikedSongs(&c.Client)
			playlist.PlayFolder(&c.Client, songs...)
		case "addliked":
			songs := playlist.GetLikedSongs(&c.Client)
			playlist.AddFolder(&c.Client, request.Data.Song, songs...)
		case "playmost":
			songs := playlist.GetMostPlayed(&c.Client)
			playlist.PlayFolder(&c.Client, songs...)
		case "addmost":
			songs := playlist.GetMostPlayed(&c.Client)
			playlist.AddFolder(&c.Client, request.Data.Song, songs...)
		case "list":
			songs := playlist.ListSongsIn(&c.Client, request.Data.Playlist)
			err = json.NewEncoder(w).Encode(songs)
			config.Log(err)
		}
		c.Client.Close()
		return
	}
	playlist := playlist.ListStoredPlaylist(c.Client)
	err = json.NewEncoder(w).Encode(playlist)
	config.Log(err)
}
