package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playback"
	"github.com/smoka7/ava/src/song"
)

type Mpd struct {
	Client config.Connection
}

type CurrentSongResponse struct {
	Album  string
	Artist string
	Date   string
	Genre  string
	Id     string
	Pos    string
	Title  string
	File   string
	Liked  bool
}

type Status struct {
	Consume    bool   `json:"consume,omitempty"`
	Duration   string `json:"duration,omitempty"`
	Elapsed    string `json:"elapsed,omitempty"`
	Mixrampdb  string `json:"mixrampdb,omitempty"`
	Random     bool   `json:"random,omitempty"`
	Repeat     bool   `json:"repeat,omitempty"`
	Single     bool   `json:"single,omitempty"`
	State      string `json:"state,omitempty"`
	Volume     string `json:"volume,omitempty"`
	Xfade      string `json:"xfade,omitempty"`
	UpdatingDB string `json:"updating_db,omitempty"`
}

type StatusResponse struct {
	Status      Status
	CurrentSong CurrentSongResponse
}

type PlaybackRequest struct {
	Command string `json:"command"`
	Data    Data   `json:"data"`
}

type Data struct {
	Start int `json:"start"`
}

type ErorrResponse struct {
	Error string `json:"error"`
}

var err error

func NewClient(c config.Connection) (cl Mpd) {
	cl.Client = c
	return
}

func (c Mpd) Playback(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var request PlaybackRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)

		err = playback.Command(c.Client, request.Command, request.Data.Start)
		sendErrorResponse(w, err)
	}
}

// writes the current mpd server status
func (c Mpd) Status(w http.ResponseWriter, r *http.Request) {
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

func newCurrentSong(song song.Song) CurrentSongResponse {
	if song.Info["Title"] == "" {
		split := strings.Split(song.Info["file"], "/")
		song.Info["Title"] = split[len(split)-1]
	}

	return CurrentSongResponse{
		Album:  song.Info["Album"],
		Artist: song.Info["Artist"],
		Date:   song.Info["Date"],
		Genre:  song.Info["Genre"],
		Id:     song.Info["Id"],
		Pos:    song.Info["Pos"],
		Title:  song.Info["Title"],
		File:   song.Info["file"],
		Liked:  song.Liked,
	}
}

func newStatus(status mpd.Attrs) Status {
	return Status{
		Consume:    status["consume"] == "1",
		Duration:   status["duration"],
		Elapsed:    status["elapsed"],
		Mixrampdb:  status["mixrampdb"],
		Random:     status["random"] == "1",
		Repeat:     status["repeat"] == "1",
		Single:     status["single"] == "1",
		State:      status["state"],
		Volume:     status["volume"],
		Xfade:      status["xfade"],
		UpdatingDB: status["updating_db"],
	}
}

// checks if the request is valid and returns an error if not
func sendErrorResponse(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(ErorrResponse{err.Error()})
		config.Log(err)
	}
}
