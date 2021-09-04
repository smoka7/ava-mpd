package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playback"
	"github.com/smoka7/ava/src/song"
)

type Mpd struct {
	Client config.Connection
}

type Request struct {
	Command string `json:"command"`
	Data    Data   `json:"data"`
}

type Data struct {
	Start    int    `json:"start"`
	Finish   int    `json:"finish"`
	Song     string `json:"song"`
	Playlist string `json:"playlist"`
}

type Response struct {
	Status      map[string]string
	CurrentSong map[string]string
	AlbumArt    string
}

var request Request

func (c *Mpd) Playback(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c.Client.Connect()
		json.NewDecoder(r.Body).Decode(&request)
		switch request.Command {
		case "next":
			playback.NextSong(c.Client)
		case "previous":
			playback.PrevSong(c.Client)
		case "toggle":
			playback.Toggle(c.Client)
		case "stop":
			playback.Stop(c.Client)
		case "single":
			playback.Single(c.Client)
		case "repeat":
			playback.Repeat(c.Client)
		case "random":
			playback.Random(c.Client)
		case "consume":
			playback.Consume(c.Client)
		case "clear":
			playback.ClearQueue(c.Client)
		case "seek":
			playback.Seek(c.Client, request.Data.Start)
		case "changeVolume":
			playback.ChangeVolume(c.Client, request.Data.Start)
		}
		c.Client.Close()
		return
	}
}

//writes the current mpd server status
func (c *Mpd) Status(w http.ResponseWriter, r *http.Request) {
	currentSong := song.GetCurrentSong(c.Client)
	response := Response{
		Status:      song.GetStatus(c.Client),
		AlbumArt:    song.ServeAlbumArt(c.Client, currentSong["file"]),
		CurrentSong: currentSong,
	}
	json.NewEncoder(w).Encode(response)
}

func (c *Mpd) Song(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&request)
		switch request.Command {
		case "albumart":
			albumArtUrl := song.ServeAlbumArt(c.Client, request.Data.Song)
			json.NewEncoder(w).Encode(
				map[string]string{
					"AlbumArt": albumArtUrl,
				})
			return
		case "like":
			song.ToggleLike(c.Client, request.Data.Song)
		}
	}
}

