package controller

import (
	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
	"github.com/smoka7/ava/src/song"
)

type Mpd struct {
	Client           config.Connection
	DownloadCoverArt bool
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
	Consume    string `json:"consume,omitempty"`
	Duration   string `json:"duration,omitempty"`
	Elapsed    string `json:"elapsed,omitempty"`
	Mixrampdb  string `json:"mixrampdb,omitempty"`
	Random     string `json:"random,omitempty"`
	Repeat     string `json:"repeat,omitempty"`
	Single     string `json:"single,omitempty"`
	State      string `json:"state,omitempty"`
	Volume     string `json:"volume,omitempty"`
	Xfade      string `json:"xfade,omitempty"`
	UpdatingDB string `json:"updating_db,omitempty"`
}

type StatusResponse struct {
	Status      Status
	CurrentSong CurrentSongResponse
}

type SearchRequest struct {
	Terms   []string `json:"terms"`
	Command string   `json:"command"`
}

type SearchResponse struct {
	Songs playlist.SearchResult
	Error error
}

var err error

func newCurrentSong(song song.Song) CurrentSongResponse {
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
		Consume:    status["consume"],
		Duration:   status["duration"],
		Elapsed:    status["elapsed"],
		Mixrampdb:  status["mixrampdb"],
		Random:     status["random"],
		Repeat:     status["repeat"],
		Single:     status["single"],
		State:      status["state"],
		Volume:     status["volume"],
		Xfade:      status["xfade"],
		UpdatingDB: status["updating_db"],
	}
}
