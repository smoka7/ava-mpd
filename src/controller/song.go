package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
	"github.com/smoka7/ava/src/song"
)

type SongRequest struct {
	Command string   `json:"command"`
	Data    SongData `json:"data"`
}

type SongData struct {
	ID   int
	File string
}

type SongInfoResponse struct {
	Info     mpd.Attrs
	Stickers []mpd.Sticker
}

type CoverArtResponse struct {
	Url string
}

func (c *Mpd) Song(w http.ResponseWriter, r *http.Request) {
	var request SongRequest
	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		c.Client.Connect()
		switch request.Command {
		case "like":
			song.ToggleLike(&c.Client, request.Data.File)
		case "info":
			response := c.getSongResponse(request.Data.ID)
			err := json.NewEncoder(w).Encode(response)
			config.Log(err)
		case "albumArt":
			response := c.getSongCover(request.Data.ID)
			err := json.NewEncoder(w).Encode(response)
			config.Log(err)
		}
		c.Client.Close()
	}
}

func (c *Mpd) getSongResponse(songPos int) SongInfoResponse {
	file, err := playlist.GetSongFile(&c.Client, songPos)
	if err != nil {
		return SongInfoResponse{}
	}

	return SongInfoResponse{
		Info:     song.NewSong().GetSongInfo(&c.Client, file).Info,
		Stickers: song.GetStickers(&c.Client, file),
	}
}

func (c *Mpd) getSongCover(songPos int) CoverArtResponse {
	file, err := playlist.GetSongFile(&c.Client, songPos)
	if err != nil {
		config.Log(err)
		return CoverArtResponse{
			Url: "default",
		}
	}

	return CoverArtResponse{
		Url: song.ServeAlbumArt(c.Client, file),
	}
}
