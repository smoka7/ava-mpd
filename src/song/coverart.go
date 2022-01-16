package song

import (
	"encoding/base32"
	"errors"
	"os"

	"github.com/smoka7/ava/src/config"
)

//returns the URL of the song
func ServeAlbumArt(c config.Connection, songPath string) (coverUrl string) {
	var song Song
	if songPath == "" {
		return "default"
	}
	c.Connect()
	defer c.Close()
	song.GetSongInfo(&c, songPath)
	coverPath, coverUrl := song.getCoverArtPath()
	// if cover isn't cached get it from mpd
	if _, err := os.Stat(coverPath); os.IsNotExist(err) {
		err := song.getAlbumArt(&c)
		if err != nil {
			config.Log(err)
			return "default"
		}
	}
	return
}

//gets the cover from mpd
func (s *Song) getAlbumArt(c *config.Connection) (err error) {
	coverBin, err := c.Client.AlbumArt(s.Info["file"])
	if err == nil {
		s.writeCoverToFile(coverBin)
		return
	}
	// get the embed cover
	coverBin, err = c.Client.ReadPicture(s.Info["file"])
	if err == nil {
		s.writeCoverToFile(coverBin)
		return
	}
	err = errors.New("cover not found")
	return
}

//writes the cover data to cache folder
func (s *Song) writeCoverToFile(data []byte) {
	coverPath, _ := s.getCoverArtPath()
	coverFile, err := os.OpenFile(coverPath, os.O_CREATE|os.O_WRONLY, 0600)
	config.Log(err)
	defer coverFile.Close()
	_, err = coverFile.Write(data)
	config.Log(err)
}

//generates a name for cover and returns its path in URL
func (s *Song) getCoverArtPath() (coverPath string, url string) {
	coverFolder, _ := os.UserCacheDir()
	coverFolder += "/ava-mpd/coverart/"
	if _, e := os.Stat(coverFolder); os.IsNotExist(e) {
		err := os.MkdirAll(coverFolder, 0777)
		config.Log(err)
	}
	fileName := sanitize(s.Info["Album"] + s.Info["AlbumArtist"])
	coverPath = coverFolder + fileName
	url = "/coverart/" + fileName
	return
}

//return a valid file name
func sanitize(name string) string {
	return base32.StdEncoding.EncodeToString([]byte(name))
}
