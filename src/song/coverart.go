package song

import (
	"encoding/base32"
	"errors"
	"os"

	"github.com/smoka7/ava/src/brainz"
	"github.com/smoka7/ava/src/config"
)

// message to send to client when we can't find any cover
const DefaultCover = "default"

// path of coverarts cache directory
var coversDirectory string

// makes the coverArt cache directory for firstTime Use
func init() {
	coversDirectory, _ = os.UserCacheDir()
	coversDirectory += "/ava-mpd/coverart/"
	if _, e := os.Stat(coversDirectory); os.IsNotExist(e) {
		err = os.MkdirAll(coversDirectory, 0o777)
		config.Log(err)
	}
}

// returns the URL of the song
func ServeAlbumArt(c config.Connection, songPath string) string {
	if songPath == "" {
		return DefaultCover
	}

	c.Connect()
	defer c.Close()

	song := NewSong()
	song.GetSongInfo(c, songPath)
	err = song.setCoverArtPath()
	if err != nil {
		return DefaultCover
	}

	//  if cover isn't cached get it from mpd
	if _, err := os.Stat(song.CoverArt.path); os.IsNotExist(err) {
		err := song.getAlbumArt(c)
		if err != nil {
			config.Log(err)
			return DefaultCover
		}
	}
	return song.CoverArt.url
}

// gets the cover from mpd
func (s *Song) getAlbumArt(c config.Connection) error {
	coverBin, err := c.Client.AlbumArt(s.Info["file"])
	if err == nil {
		s.writeCoverToFile(coverBin)
		return nil
	}
	//  get the embed cover
	coverBin, err = c.Client.ReadPicture(s.Info["file"])
	if err == nil {
		s.writeCoverToFile(coverBin)
		return nil
	}

	if releaseID, ok := s.Info["MUSICBRAINZ_ALBUMID"]; ok &&
		releaseID != "" &&
		c.DownloadCoverArt {
		coverBin, err = brainz.GetCover(releaseID)
		if err == nil {
			s.writeCoverToFile(coverBin)
			return nil
		}
	}

	err = errors.New("cover not found")
	return err
}

// writes the cover data to cache folder
func (s *Song) writeCoverToFile(data []byte) {
	err = os.WriteFile(s.CoverArt.path, data, 0o600)
	config.Log(err)
}

// finds the path and URL of the cover art
func (s *Song) setCoverArtPath() error {
	fileName := sanitize(s.Info["Album"] + s.Info["AlbumArtist"])

	if fileName == "" {
		fileName = sanitize(s.Info["Title"] + s.Info["Artist"])
	}

	if fileName == "" {
		err = errors.New("couldn't find a name for cover")
		config.Log(err)
		return err
	}

	s.CoverArt.path = coversDirectory + fileName
	s.CoverArt.url = "/coverart/" + fileName
	return nil
}

// return a valid file name
func sanitize(name string) string {
	return base32.StdEncoding.EncodeToString([]byte(name))
}
