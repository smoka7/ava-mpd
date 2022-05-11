package song

import (
	"encoding/base32"
	"errors"
	"os"

	"github.com/smoka7/ava/src/brainz"
	"github.com/smoka7/ava/src/config"
)

const DefaultCover = "default"

// returns the URL of the song
func ServeAlbumArt(c config.Connection, songPath string) (coverURL string) {
	song := NewSong()
	if songPath == "" {
		return DefaultCover
	}
	c.Connect()
	defer c.Close()

	song.GetSongInfo(&c, songPath)
	coverPath, coverURL, err := song.getCoverArtPath()
	if err != nil {
		return DefaultCover
	}

	//  if cover isn't cached get it from mpd
	if _, err := os.Stat(coverPath); os.IsNotExist(err) {
		err := song.getAlbumArt(&c)
		if err != nil {
			config.Log(err)
			return DefaultCover
		}
	}
	return
}

// gets the cover from mpd
func (s *Song) getAlbumArt(c *config.Connection) (err error) {
	coverBin, err := c.Client.AlbumArt(s.Info["file"])
	if err == nil {
		s.writeCoverToFile(coverBin)
		return
	}
	//  get the embed cover
	coverBin, err = c.Client.ReadPicture(s.Info["file"])
	if err == nil {
		s.writeCoverToFile(coverBin)
		return
	}

	if releaseID, ok := s.Info["MUSICBRAINZ_ALBUMID"]; ok &&
		releaseID != "" &&
		c.DownloadCoverArt {
		coverBin, err = brainz.GetCover(releaseID)
		if err == nil {
			s.writeCoverToFile(coverBin)
			return
		}
	}

	err = errors.New("cover not found")
	return
}

// writes the cover data to cache folder
func (s *Song) writeCoverToFile(data []byte) {
	coverPath, _, err := s.getCoverArtPath()
	if err != nil {
		return
	}

	coverFile, err := os.OpenFile(coverPath, os.O_CREATE|os.O_WRONLY, 0600)
	config.Log(err)
	defer coverFile.Close()
	_, err = coverFile.Write(data)
	config.Log(err)
}

// generates a name for cover and returns its path in URL
func (s *Song) getCoverArtPath() (coverPath, url string, err error) {
	coverFolder, _ := os.UserCacheDir()
	coverFolder += "/ava-mpd/coverart/"
	if _, e := os.Stat(coverFolder); os.IsNotExist(e) {
		err := os.MkdirAll(coverFolder, 0777)
		config.Log(err)
	}
	fileName := sanitize(s.Info["Album"] + s.Info["AlbumArtist"])

	if fileName == "" {
		fileName = sanitize(s.Info["Title"] + s.Info["Artist"])
	}

	if fileName == "" {
		err = errors.New("couldn't find a name for cover")
		config.Log(err)
		return
	}

	coverPath = coverFolder + fileName
	url = "/coverart/" + fileName
	return
}

// return a valid file name
func sanitize(name string) string {
	return base32.StdEncoding.EncodeToString([]byte(name))
}
