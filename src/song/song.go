package song

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

type CoverArt struct {
	url, path string
}

type Song struct {
	Info     SongInfo
	CoverArt CoverArt
	Liked    bool
}
type SongInfo struct {
	Id                         string `json:"Id,omitempty"`
	Pos                        string `json:"Pos,omitempty"`
	Album                      string `json:"album,omitempty"`
	Albumartist                string `json:"albumartist,omitempty"`
	Albumartistsort            string `json:"albumartistsort,omitempty"`
	Albumsort                  string `json:"albumsort,omitempty"`
	Artist                     string `json:"artist,omitempty"`
	Artistsort                 string `json:"artistsor,omitempty"`
	Comment                    string `json:"comment,omitempty"`
	Composer                   string `json:"composer,omitempty"`
	Composersort               string `json:"composersort,omitempty"`
	Conductor                  string `json:"conductor,omitempty"`
	Date                       string `json:"date,omitempty"`
	Disc                       string `json:"disc,omitempty"`
	Ensemble                   string `json:"ensemble,omitempty"`
	File                       string `json:"file,omitempty"`
	Genre                      string `json:"genre,omitempty"`
	Grouping                   string `json:"grouping,omitempty"`
	Label                      string `json:"label,omitempty"`
	Location                   string `json:"location,omitempty"`
	Mood                       string `json:"mood,omitempty"`
	Movement                   string `json:"movement,omitempty"`
	Movementnumber             string `json:"movementnumber,omitempty"`
	Musicbrainz_albumartistid  string `json:"musicbrainz_albumartistid,omitempty"`
	Musicbrainz_albumid        string `json:"musicbrainz_albumid,omitempty"`
	Musicbrainz_artistid       string `json:"musicbrainz_artistid,omitempty"`
	Musicbrainz_releasetrackid string `json:"musicbrainz_releasetracki,omitempty"`
	Musicbrainz_trackid        string `json:"musicbrainz_trackid,omitempty"`
	Musicbrainz_workid         string `json:"musicbrainz_workid,omitempty"`
	Name                       string `json:"name,omitempty"`
	Originaldate               string `json:"originaldate,omitempty"`
	Performer                  string `json:"performer,omitempty"`
	Title                      string `json:"title,omitempty"`
	Titlesort                  string `json:"titlesort,omitempty"`
	Track                      string `json:"track,omitempty"`
	Work                       string `json:"work,omitempty"`
}

func newInfo(i mpd.Attrs) SongInfo {
	return SongInfo{
		Id:                         i["Id"],
		Pos:                        i["Pos"],
		Album:                      i["Album"],
		Albumartist:                i["Albumartist"],
		Albumartistsort:            i["Albumartistsort"],
		Albumsort:                  i["Albumsort"],
		Artist:                     i["Artist"],
		Artistsort:                 i["Artistsort"],
		Comment:                    i["Comment"],
		Composer:                   i["Composer"],
		Composersort:               i["Composersort"],
		Conductor:                  i["Conductor"],
		Date:                       i["Date"],
		Disc:                       i["Disc"],
		Ensemble:                   i["Ensemble"],
		File:                       i["file"],
		Genre:                      i["Genre"],
		Grouping:                   i["Grouping"],
		Label:                      i["Label"],
		Location:                   i["Location"],
		Mood:                       i["Mood"],
		Movement:                   i["Movement"],
		Movementnumber:             i["Movementnumber"],
		Musicbrainz_albumartistid:  i["MUSICBRAINZ_ALBUMARTISTID"],
		Musicbrainz_albumid:        i["MUSICBRAINZ_ALBUMID"],
		Musicbrainz_artistid:       i["MUSICBRAINZ_ARTISTID"],
		Musicbrainz_releasetrackid: i["MUSICBRAINZ_RELEASETRACKId"],
		Musicbrainz_trackid:        i["MUSICBRAINZ_TRACKID"],
		Musicbrainz_workid:         i["MUSICBRAINZ_WORKID"],
		Name:                       i["Name"],
		Originaldate:               i["Originaldate"],
		Performer:                  i["Performer"],
		Title:                      i["Title"],
		Titlesort:                  i["Titlesort"],
		Track:                      i["Track"],
		Work:                       i["Work"],
	}
}

func NewSong(c config.Connection, path string) Song {
	song := Song{}
	song.setSongInfo(c, path)
	return song
}

// returns the server status.
func GetStatus(c config.Connection) (status mpd.Attrs) {
	c.Connect()
	defer c.Close()
	status, err := c.Status()
	config.Log(err)
	return
}

// returns the currentsong info.
func GetCurrentSong(c config.Connection) (song Song) {
	c.Connect()
	defer c.Close()

	info, err := c.CurrentSong()
	config.Log(err)
	song.Info = newInfo(info)

	liked := getSticker(c, song.Info.File, "liked")
	if liked != nil && liked.Value != "false" {
		song.Liked = true
		return
	}
	song.Liked = false

	return
}

// gets the song info.
func (s *Song) setSongInfo(c config.Connection, file string) {
	info, err := c.ListAllInfo(file)
	config.Log(err)
	if len(info) > 0 {
		s.Info = newInfo(info[0])
	}
}

// set the sticker name to value for song.
func setSticker(c config.Connection, song, name, value string) {
	err := c.StickerSet(song, name, value)
	config.Log(err)
}

// returns all the stickers of a song.
func GetStickers(c config.Connection, file string) []mpd.Sticker {
	stickers, err := c.StickerList(file)
	config.Log(err)
	return stickers
}

// return the sticker name of a song.
func getSticker(c config.Connection, file, name string) (status *mpd.Sticker) {
	status, err := c.StickerGet(file, name)
	if err != nil && err.Error() != "command 'sticker' failed: no such sticker" {
		config.Log(err)
	}
	return
}

// Increments played count of current song.
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

// sets last played time of song.
func setLastPlayed(c config.Connection, uri string) {
	now := time.Now().Unix()
	setSticker(c, uri, "lastplayed", fmt.Sprintf("%d", now))
}

// toggles like state of song.
func ToggleLike(c config.Connection, uri string) {
	// when uri is empty toggle current song
	if uri == "" {
		uri = GetCurrentSong(c).Info.File
	}

	liked := getSticker(c, uri, "liked")
	if liked != nil && liked.Value == "true" {
		setSticker(c, uri, "liked", "false")
		return
	}
	setSticker(c, uri, "liked", "true")
}
