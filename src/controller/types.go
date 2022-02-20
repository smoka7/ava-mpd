package controller

import (
	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
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

type StatusResponse struct {
	Status      mpd.Attrs
	CurrentSong song.Song
	AlbumArt    string
}

type SearchRequest struct {
	Terms   []string `json:"terms"`
	Command string   `json:"command"`
}

type SearchResponse struct {
	Songs playlist.SearchResult
	Error error
}

type SongInfoResponse struct {
	Info     mpd.Attrs
	Stickers []mpd.Sticker
	AlbumArt string
}

type SettingsResponse struct {
	Outputs       []mpd.Attrs
	DatabaseStats mpd.Attrs
	ReplayGain    string
}

var err error
