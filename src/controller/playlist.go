package controller

import (
	"encoding/json"
	"errors"
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
			playlist.AddFolder(&c.Client, songs...)
		case "playmost":
			songs := playlist.GetMostPlayed(&c.Client)
			playlist.PlayFolder(&c.Client, songs...)
		case "addmost":
			songs := playlist.GetMostPlayed(&c.Client)
			playlist.AddFolder(&c.Client, songs...)
		case "list":
			songs := playlist.ListSongs(&c.Client, request.Data.Playlist)
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

func (c *Mpd) Queue(w http.ResponseWriter, r *http.Request) {
	c.Client.Connect()
	defer c.Client.Close()
	if r.Method == http.MethodPost {
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "addsong":
			file := playlist.GetSongFile(&c.Client, request.Data.Start)
			playlist.AddSongToPlaylist(&c.Client, request.Data.Playlist, file)
		case "delete":
			playlist.DeleteSong(&c.Client, request.Data.Start, -1)
		case "move":
			playlist.MoveSong(&c.Client, request.Data.Start, request.Data.Finish)
		case "play":
			playlist.PlaySong(&c.Client, request.Data.Start)
		case "removeduplicate":
			playlist.RemoveDuplicatesongs(&c.Client, "")
		case "save":
			playlist.SaveQueue(&c.Client, request.Data.Playlist)
		}
		return
	}
	playlist := playlist.GetQueue(c.Client)
	err = json.NewEncoder(w).Encode(playlist)
	config.Log(err)
}

func (c *Mpd) ServerFolders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.Client.Connect()
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "add":
			playlist.AddFolder(&c.Client, request.Data.Playlist)
		case "play":
			playlist.PlayFolder(&c.Client, request.Data.Playlist)
		case "list":
			folders := playlist.ListFolders(&c.Client, request.Data.Playlist)
			err = json.NewEncoder(w).Encode(folders)
			config.Log(err)
		}
		c.Client.Close()
		return
	}
}

func (c *Mpd) SearchServer(w http.ResponseWriter, r *http.Request) {
	c.Client.Connect()
	defer c.Client.Close()
	if r.Method == http.MethodPost {
		var request SearchRequest
		var files playlist.SearchResult
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "search":
			files, err = playlist.SearchServer(&c.Client, request.Terms...)
		case "searchadd":
			err = playlist.SearchAdd(&c.Client, request.Terms...)
		case "searchplay":
			err = playlist.SearchPlay(&c.Client, request.Terms...)
		default:
			err = errors.New("chose a command")
		}
		err = json.NewEncoder(w).Encode(SearchResponse{
			Songs: files,
			Error: err,
		})
		config.Log(err)
		return
	}
}
