package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playback"
	"github.com/smoka7/ava/src/playlist"
	"github.com/smoka7/ava/src/song"
)

var request Request

func NewClient() (cl Mpd) {
	cl.Client.ReadEnv()
	cl.Client.ReadFromFlags()
	if cl.Client.Address == "" {
		log.Println("mpd server address is empty")
		os.Exit(1)
	}
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
		Status:      song.GetStatus(c.Client),
		AlbumArt:    song.ServeAlbumArt(c.Client, currentSong.Info["file"]),
		CurrentSong: currentSong,
	}
	c.Client.Close()
	err := json.NewEncoder(w).Encode(response)
	config.Log(err)
}

func (c *Mpd) Song(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		c.Client.Connect()
		switch request.Command {
		case "like":
			song.ToggleLike(&c.Client, request.Data.Song)
		case "info":
			response := c.getSongResponse(request.Data.Start)
			err := json.NewEncoder(w).Encode(response)
			config.Log(err)
		}
		c.Client.Close()
	}
}

func (c *Mpd) Settings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.Client.Connect()
		defer c.Client.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "crossfade":
			c.Client.ChangeCrossfade(request.Data.Start)
		case "download":
			c.Client.ToggleDownloadCover()
		case "mixrampdb":
			c.Client.ChangeMixRampdb(request.Data.Start)
		case "enableOutput":
			c.Client.EnableOutput(request.Data.Start)
		case "disableOutput":
			c.Client.DisableOutput(request.Data.Start)
		case "deleteCache":
			config.DeleteCache()
		case "setGain":
			c.Client.ChangeReplayGain(request.Data.Start)
		case "updateDatabase":
			c.Client.UpdateDatabase()
		}
		return
	}
	response := SettingsResponse{
		DatabaseStats:    c.Client.DatabaseStats(),
		ReplayGain:       c.Client.GetReplayGain(),
		Outputs:          c.Client.ListOutputs(),
		DownloadCoverArt: c.Client.DownloadCoverArt,
	}
	err := json.NewEncoder(w).Encode(response)
	config.Log(err)
}

func (c *Mpd) getSongResponse(songPos int) SongInfoResponse {
	file, err := playlist.GetSongFile(&c.Client, songPos)
	if err != nil {
		return SongInfoResponse{}
	}

	return SongInfoResponse{
		Info:     song.NewSong().GetSongInfo(&c.Client, file).Info,
		Stickers: song.GetStickers(&c.Client, file),
		AlbumArt: song.ServeAlbumArt(c.Client, file),
	}
}
