package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
)

type QueueRequest struct {
	Command string    `json:"command"`
	Data    QueueData `json:"data"`
}

type QueueData struct {
	ID            int
	Start, Finish int
	Playlist      string
}

func (c *Mpd) Queue(w http.ResponseWriter, r *http.Request) {
	var request QueueRequest
	c.Client.Connect()
	defer c.Client.Close()
	if r.Method == http.MethodPost {
		err = json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "addsong":
			file, err := playlist.GetSongFile(&c.Client, request.Data.ID)
			if err != nil {
				return
			}
			playlist.AddSongToPlaylist(&c.Client, request.Data.Playlist, file)
		case "delete":
			playlist.DeleteSong(&c.Client, request.Data.ID)
		case "move":
			playlist.MoveSong(&c.Client, request.Data.Start, request.Data.Finish)
		case "play":
			playlist.PlaySong(&c.Client, request.Data.ID)
		case "shuffleAlbum":
			playlist.ShuffleAlbum(&c.Client, request.Data.Start)
		case "removeduplicate":
			playlist.RemoveDuplicatesongs(&c.Client, "")
		case "save":
			playlist.SaveQueue(&c.Client, request.Data.Playlist)
		}
		return
	}
	err = r.ParseForm()
	config.Log(err)

	playlist := playlist.GetQueue(c.Client, r.Form.Get("page"))
	err = json.NewEncoder(w).Encode(playlist)
	config.Log(err)
}
