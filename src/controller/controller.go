package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playback"
	"github.com/smoka7/ava/src/song"
)

var request Request

func NewClient() (cl Mpd) {
	cl.Client.ReadConfigs()
	return
}

func (c *Mpd) Playback(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		playback.Command(c.Client, request.Command, request.Data.Start)
	}
}

// writes the current mpd server status
func (c *Mpd) Status(w http.ResponseWriter, r *http.Request) {
	currentSong := song.GetCurrentSong(c.Client)
	c.Client.Connect()
	response := StatusResponse{
		Status:      newStatus(song.GetStatus(c.Client)),
		CurrentSong: newCurrentSong(currentSong),
	}
	c.Client.Close()
	err := json.NewEncoder(w).Encode(response)
	config.Log(err)
}
