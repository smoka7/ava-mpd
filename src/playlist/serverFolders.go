package playlist

import (
	"fmt"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

type FolderRequest struct {
	Command string     `json:"command"`
	Data    FolderData `json:"data"`
}

type FolderData struct {
	Pos  string
	File string
}

// returns content of the folder
func (a action) ListFolders(d FolderData) (list ServerList) {
	contents, err := a.ListInfo(d.File)
	config.Log(err)
	list.Files = make(Files, 0)
	list.Directories = make(Directories, 0)
	for _, item := range contents {
		// ignore playlists
		if _, ok := item["playlist"]; ok {
			continue
		}

		if _, ok := item["file"]; ok {
			list.Files = append(list.Files, newFile(item))
			continue
		}
		list.Directories = append(list.Directories, newFolder(item))
	}
	return
}

// clears the current queue and plays the folder
func (a action) PlayFolder(uris ...string) {
	cm := a.BeginCommandList()
	cm.Clear()
	for _, uri := range uris {
		cm.Add(uri)
	}
	cm.Play(0)
	err := cm.End()
	config.Log(err)
}

// adds the folder to the current queue
func (a action) AddFolder(pos string, uris ...string) {
	add := func(pos string) {
		for _, uri := range uris {
			err := a.Command("add %s %s", uri, pos).OK()
			config.Log(err)
		}
	}
	current := a.getCurrentSongPos()
	if current < 0 {
		a.addToEndOfQueue(uris)
		return
	}

	switch pos {
	case "currentSong":
		add("+0")

	case "endOfQueue":
		a.addToEndOfQueue(uris)

	case "currentAlbum":
		_, endOfAlbum := a.findSongsAlbum(current)
		add(fmt.Sprint(endOfAlbum))
	}
}

func (a action) addToEndOfQueue(uris []string) {
	list := a.BeginCommandList()
	for _, uri := range uris {
		list.Add(uri)
	}
	err := list.End()
	config.Log(err)
}

func newFolder(item mpd.Attrs) Directory {
	return Directory{
		Directory: item["directory"],
	}
}

func newFile(item mpd.Attrs) File {
	return File{
		File: item["file"],
	}
}
