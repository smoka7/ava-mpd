package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
	"github.com/smoka7/ava/src/song"
)

type songCommand string

const (
	like     songCommand = "like"
	info     songCommand = "info"
	albumArt songCommand = "albumArt"
)

type SongRequest struct {
	Command songCommand `json:"command"`
	Data    SongData    `json:"data"`
}

type SongData struct {
	File string `json:""`
	ID   int    `json:""`
}

type SongInfoResponse struct {
	Info     song.SongInfo `json:""`
	Stickers []mpd.Sticker `json:""`
}

type CoverArtResponse struct {
	Url string `json:""`
}

func (c Mpd) Song(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	var request SongRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	config.Log(err)

	c.Connect()
	defer c.Close()

	switch request.Command {
	case like:
		song.ToggleLike(c.Connection, request.Data.File)
	case info:
		response := c.getSongResponse(request.Data.ID)
		err := json.NewEncoder(w).Encode(response)
		config.Log(err)
	case albumArt:
		response := c.getSongCover(request.Data.ID)
		err := json.NewEncoder(w).Encode(response)
		config.Log(err)
	}
}

func (c Mpd) getSongResponse(songPos int) SongInfoResponse {
	a := playlist.NewAction(c.Connection)
	file, err := a.GetSongFile(songPos)
	if err != nil {
		return SongInfoResponse{}
	}

	return SongInfoResponse{
		Info:     song.NewSong(c.Connection, file).Info,
		Stickers: song.GetStickers(c.Connection, file),
	}
}

func (c Mpd) getSongCover(songPos int) CoverArtResponse {
	a := playlist.NewAction(c.Connection)
	file, err := a.GetSongFile(songPos)
	if err != nil {
		config.Log(err)
		return CoverArtResponse{
			Url: "default",
		}
	}

	return CoverArtResponse{
		Url: song.ServeAlbumArt(c.Connection, file),
	}
}
