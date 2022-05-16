package song

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

var err error

type Song struct {
	Info  mpd.Attrs
	Liked bool
}

func NewSong() *Song {
	return &Song{}
}

// returns the server status
func GetStatus(c config.Connection) (status mpd.Attrs) {
	c.Connect()
	defer c.Close()
	status, err = c.Client.Status()
	config.Log(err)
	return
}

// returns the currentsong info
func GetCurrentSong(c config.Connection) (song Song) {
	c.Connect()
	defer c.Close()
	song.Info, err = c.Client.CurrentSong()
	config.Log(err)
	liked := GetSticker(&c, song.Info["file"], "liked")
	if liked != nil && liked.Value != "false" {
		song.Liked = true
		return
	}
	song.Liked = false
	return
}

// gets the song info
func (s *Song) GetSongInfo(c *config.Connection, file string) *Song {
	info, err := c.Client.ListAllInfo(file)
	config.Log(err)
	if len(info) > 0 {
		s.Info = info[0]
	}
	return s
}

// set the sticker name to value for song
func SetSticker(c *config.Connection, song, name, value string) {
	err = c.Client.StickerSet(song, name, value)
	config.Log(err)
}

// returns all the stickers of a song
func GetStickers(c *config.Connection, file string) (status []mpd.Sticker) {
	status, err = c.Client.StickerList(file)
	config.Log(err)
	return
}

// return the sticker name of a song
func GetSticker(c *config.Connection, file, name string) (status *mpd.Sticker) {
	status, err = c.Client.StickerGet(file, name)
	if err != nil && err.Error() != "command 'sticker' failed: no such sticker" {
		config.Log(err)
	}
	return
}

// Increments played count of current song
func IncrementPCount(c *config.Connection, filepath string) {
	SetLastPlayed(c, filepath)
	playedCount := GetSticker(c, filepath, "playedcount")
	if playedCount == nil {
		SetSticker(c, filepath, "playedcount", "1")
		return
	}
	if value, err := strconv.Atoi(playedCount.Value); err == nil {
		SetSticker(c, filepath, "playedcount", fmt.Sprintf("%d", value+1))
	}
}

// sets last played time of song
func SetLastPlayed(c *config.Connection, uri string) {
	now := time.Now().Unix()
	SetSticker(c, uri, "lastplayed", fmt.Sprintf("%d", now))
}

// toggles like state of song
func ToggleLike(c *config.Connection, uri string) {
	liked := GetSticker(c, uri, "liked")
	if liked != nil && liked.Value == "true" {
		SetSticker(c, uri, "liked", "false")
		return
	}
	SetSticker(c, uri, "liked", "true")
}
