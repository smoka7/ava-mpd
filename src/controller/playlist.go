package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/playlist"
)

type SearchRequest struct {
	Terms []string `json:"terms"`
}

func (c *Mpd) StoredPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&request)
		switch request.Command {
		case "load":
			playlist.LoadPlaylist(c.Client, request.Data.Playlist)
		case "clear":
			playlist.ClearPlaylist(c.Client, request.Data.Playlist)
		case "delete":
			playlist.DeletePlaylist(c.Client, request.Data.Playlist)
		case "play":
			playlist.PlayPlaylist(c.Client, request.Data.Playlist)
		case "rename":
			playlist.RenamePlaylist(c.Client, request.Data.Playlist, request.Data.Song)
		case "playliked":
			songs := playlist.GetLikedSongs(c.Client)
			playlist.PlayFolder(c.Client, songs...)
		case "addliked":
			songs := playlist.GetLikedSongs(c.Client)
			playlist.AddFolder(c.Client, songs...)
		case "playmost":
			songs := playlist.GetMostPlayed(c.Client)
			playlist.PlayFolder(c.Client, songs...)
		case "addmost":
			songs := playlist.GetMostPlayed(c.Client)
			playlist.AddFolder(c.Client, songs...)
		case "list":
			songs := playlist.ListSongs(c.Client, request.Data.Playlist)
			json.NewEncoder(w).Encode(songs)
		}
		return
	}
	playlist := playlist.ListStoredPlaylist(c.Client)
	json.NewEncoder(w).Encode(playlist)
}

func (c *Mpd) Queue(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&request)
		switch request.Command {
		case "addsong":
			playlist.AddSongToPlaylist(c.Client, request.Data.Playlist, request.Data.Song)
		case "delete":
			playlist.DeleteSong(c.Client, request.Data.Start, -1)
		case "move":
			playlist.MoveSong(c.Client, request.Data.Start, request.Data.Finish)
		case "play":
			playlist.PlaySong(c.Client, request.Data.Start)
		case "removeduplicate":
			playlist.RemoveDuplicatesongs(c.Client)
		case "save":
			playlist.SaveQueue(c.Client, request.Data.Playlist)
		}
		return
	}
	playlist := playlist.GetQueue(c.Client)
	json.NewEncoder(w).Encode(playlist)

}

func (c *Mpd) ServerFolders(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&request)
		switch request.Command {
		case "add":
			playlist.AddFolder(c.Client, request.Data.Playlist)
		case "play":
			playlist.PlayFolder(c.Client, request.Data.Playlist)
		case "list":
			folders := playlist.ListFolders(c.Client, request.Data.Playlist)
			json.NewEncoder(w).Encode(folders)
		}
		return
	}
	folders := playlist.ListFolders(c.Client, "")
	json.NewEncoder(w).Encode(folders)
}

func (c *Mpd) SearchServer(w http.ResponseWriter, r *http.Request) {
	c.Client.Connect()
	var request SearchRequest
	defer c.Client.Close()
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&request)
		files := playlist.SearchServer(c.Client, request.Terms...)
		json.NewEncoder(w).Encode(files)
		return
	}
}
