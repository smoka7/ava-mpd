package song

import (
	"errors"
	"os"

	"github.com/smoka7/ava/src/config"
)

//returns the URL of the song
func ServeAlbumArt(c config.Connection, songPath string) (coverUrl string) {
	var song Song
	song.GetSongInfo(c, songPath)
	coverPath, coverUrl := song.getCoverArtPath()
	// if cover isn't cached get it from mpd
	if _, err := os.Stat(coverPath); os.IsNotExist(err) {
		err := song.getAlbumArt(c)
		if err != nil {
			config.Log(err)
			return "default"
		}
	}
	return
}

//gets the cover from mpd
func (s *Song) getAlbumArt(c config.Connection) (err error) {
	c.Connect()
	defer c.Close()
	coverBin, err := c.Client.AlbumArt(s.Info["file"])
	if err == nil {
		s.writeCoverToFile(coverBin)
		return
	}
	// get the embed cover
	coverBin, err = s.readEmbedCover(c)
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
	coverPath = coverFolder + s.Info["Album"] + "-" + s.Info["Artist"]
	url = "/coverart/" + s.Info["Album"] + "-" + s.Info["Artist"]
	return
}

//reads the song's embed cover
func (s *Song) readEmbedCover(c config.Connection) (data []byte, err error) {
	//copied from https://github.com/fhs/gompd/blob/5f471998a4fb7c6cb9f0eba81ae1119b2eb20161/mpd/client.go#L1027
	offset := 0
	for {
		chunk, size, err := c.Client.Command("readpicture %s %d", s.Info["file"], offset).Binary()
		if err != nil {
			return data, err
		}
		data = append(data, chunk...)
		offset = len(data)
		if offset >= size {
			break
		}
	}
	return data, nil
}
