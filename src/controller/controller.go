package controller

import (
	"encoding/json"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playback"
	"github.com/smoka7/ava/src/playlist"
	"github.com/smoka7/ava/src/song"
)

var request Request

func (c *Mpd) Playback(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c.Client.Connect()
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "next":
			playback.NextSong(&c.Client)
		case "previous":
			playback.PrevSong(&c.Client)
		case "toggle":
			playback.Toggle(&c.Client)
		case "stop":
			playback.Stop(&c.Client)
		case "single":
			playback.Single(&c.Client)
		case "repeat":
			playback.Repeat(&c.Client)
		case "random":
			playback.Random(&c.Client)
		case "consume":
			playback.Consume(&c.Client)
		case "clear":
			playback.ClearQueue(&c.Client)
		case "seek":
			playback.Seek(&c.Client, request.Data.Start)
		case "changeVolume":
			playback.ChangeVolume(&c.Client, request.Data.Start)
		}
		c.Client.Close()
	}
}

//writes the current mpd server status
func (c *Mpd) Status(w http.ResponseWriter, r *http.Request) {
	currentSong := song.GetCurrentSong(c.Client)
	c.Client.Connect()
	response := StatusResponse{
		Status:      song.GetStatus(&c.Client),
		AlbumArt:    song.ServeAlbumArt(c.Client, currentSong["file"]),
		CurrentSong: currentSong,
	}
	c.Client.Close()
	err := json.NewEncoder(w).Encode(response)
	config.Log(err)
}

func (c *Mpd) Song(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		c.Client.Connect()
		switch request.Command {
		case "like":
			song.ToggleLike(&c.Client, request.Data.Song)
		case "info":
			file := playlist.GetSongFile(&c.Client, request.Data.Start)
			info := song.NewSong().GetSongInfo(&c.Client, file)
			response := SongInfoResponse{
				Info:     info.Info,
				Stickers: song.GetStickers(&c.Client, file),
				AlbumArt: song.ServeAlbumArt(c.Client, file),
			}
			err := json.NewEncoder(w).Encode(response)
			config.Log(err)
		}
		c.Client.Close()
	}
}

func (c *Mpd) Settings(w http.ResponseWriter, r *http.Request) {
	c.Client.Connect()
	defer c.Client.Close()
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "crossfade":
			config.ChangeCrossfade(c.Client, request.Data.Start)
		case "enableOutput":
			config.EnableOutput(c.Client, request.Data.Start)
		case "disableOutput":
			config.DisableOutput(c.Client, request.Data.Start)
		case "deleteCache":
			config.DeleteCache()
		case "setGain":
			config.ChangeReplayGain(c.Client, request.Data.Start)
		case "updateDatabase":
			config.UpdateDatabase(c.Client)
		}
		return
	}
	response := SettingsResponse{
		DatabaseStats: config.DatabaseStats(c.Client),
		ReplayGain:    config.GetReplayGain(c.Client),
		Outputs:       config.ListOutputs(c.Client),
	}
	err := json.NewEncoder(w).Encode(response)
	config.Log(err)
}
