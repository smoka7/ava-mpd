package song

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

var err error

type CoverArt struct {
	url, path string
}

type Song struct {
	Info     mpd.Attrs
	Liked    bool
	CoverArt CoverArt
}

func NewSong(c config.Connection, path string) Song {
	song := Song{}
	song.setSongInfo(c, path)
	return song
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
	liked := getSticker(c, song.Info["file"], "liked")
	if liked != nil && liked.Value != "false" {
		song.Liked = true
		return
	}
	song.Liked = false
	return
}

// gets the song info
func (s *Song) setSongInfo(c config.Connection, file string) {
	info, err := c.Client.ListAllInfo(file)
	config.Log(err)
	if len(info) > 0 {
		s.Info = info[0]
	}
}

// set the sticker name to value for song
func setSticker(c config.Connection, song, name, value string) {
	err = c.Client.StickerSet(song, name, value)
	config.Log(err)
}

// returns all the stickers of a song
func GetStickers(c config.Connection, file string) []mpd.Sticker {
	stickers, err := c.Client.StickerList(file)
	config.Log(err)
	return stickers
}

// return the sticker name of a song
func getSticker(c config.Connection, file, name string) (status *mpd.Sticker) {
	status, err = c.Client.StickerGet(file, name)
	if err != nil && err.Error() != "command 'sticker' failed: no such sticker" {
		config.Log(err)
	}
	return
}

// Increments played count of current song
func IncrementPCount(c config.Connection, filepath string) {
	setLastPlayed(c, filepath)
	playedCount := getSticker(c, filepath, "playedcount")
	if playedCount == nil {
		setSticker(c, filepath, "playedcount", "1")
		return
	}
	if value, err := strconv.Atoi(playedCount.Value); err == nil {
		setSticker(c, filepath, "playedcount", fmt.Sprintf("%d", value+1))
	}
}

// sets last played time of song
func setLastPlayed(c config.Connection, uri string) {
	now := time.Now().Unix()
	setSticker(c, uri, "lastplayed", fmt.Sprintf("%d", now))
}

// toggles like state of song
func ToggleLike(c config.Connection, uri string) {
	// when uri is empty toggle current song
	if uri == "" {
		uri = GetCurrentSong(c).Info["file"]
	}

	liked := getSticker(c, uri, "liked")
	if liked != nil && liked.Value == "true" {
		setSticker(c, uri, "liked", "false")
		return
	}
	setSticker(c, uri, "liked", "true")
}
