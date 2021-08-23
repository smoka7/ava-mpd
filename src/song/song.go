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
	Info mpd.Attrs
}

//returns the server status
func GetStatus(c config.Connection) (status map[string]string) {
	c.Connect()
	status, err := c.Client.Status()
	config.Log(err)
	c.Close()
	return
}

//returns the currentsong info
func GetCurrentSong(c config.Connection) (status map[string]string) {
	c.Connect()
	status, err := c.Client.CurrentSong()
	c.Close()
	config.Log(err)
	liked := GetSticker(c, status["file"], "liked")
	if liked != nil {
		status["liked"] = liked.Value
		return
	}
	status["liked"] = "false"
	return
}

//gets the song info
func (s *Song) GetSongInfo(c config.Connection, file string) {
	c.Connect()
	info, err := c.Client.ListAllInfo(file)
	config.Log(err)
	c.Close()
	if len(info) > 0 {
		s.Info = info[0]
	}
}

//set the sticker name to value for song
func SetSticker(c config.Connection, song, name, value string) {
	c.Connect()
	err := c.Client.StickerSet(song, name, value)
	config.Log(err)
	c.Close()
}

//returns all the stickers of a song
func GetStickers(c config.Connection, file string) (status []mpd.Sticker) {
	c.Connect()
	status, err := c.Client.StickerList(file)
	config.Log(err)
	c.Close()
	return
}

//return the sticker name of a song
func GetSticker(c config.Connection, file, name string) (status *mpd.Sticker) {
	c.Connect()
	status, err := c.Client.StickerGet(file, name)
	config.Log(err)
	c.Close()
	return
}

//Increments played count of current song
func IncrementPCount(c config.Connection) {
	song := GetCurrentSong(c)
	SetLastPlayed(c, song["file"])
	playedCount := GetSticker(c, song["file"], "playedcount")
	if playedCount == nil {
		SetSticker(c, song["file"], "playedcount", "1")
		return
	}
	if value, err := strconv.Atoi(playedCount.Value); err == nil {
		SetSticker(c, song["file"], "playedcount", fmt.Sprintf("%d", value+1))
	}
}

//sets last played time of song
func SetLastPlayed(c config.Connection, uri string) {
	now := time.Now().Unix()
	SetSticker(c, uri, "lastplayed", fmt.Sprintf("%d", now))
}

//toggles like state of song
func ToggleLike(c config.Connection, uri string) {
	liked := GetSticker(c, uri, "liked")
	if liked != nil && liked.Value == "true" {
		SetSticker(c, uri, "liked", "false")
		return
	}
	SetSticker(c, uri, "liked", "true")
}
