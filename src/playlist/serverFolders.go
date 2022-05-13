package playlist

import (
	"fmt"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

// returns content of the folder
func ListFolders(c *config.Connection, folder string) (list ServerList) {
	contents, err := c.Client.ListInfo(folder)
	config.Log(err)
	list.Files = make(Files, 0)
	list.Folders = make(Folders, 0)
	for _, item := range contents {
		// ignore playlists
		if _, ok := item["playlist"]; ok {
			continue
		}
		// append files to the end of list
		if _, ok := item["file"]; ok {
			list.Files = append(list.Files, newFile(item))
			continue
		}
		list.Folders = append(list.Folders, newFolder(item))
	}
	return
}

// clears the current queue and plays the folder
func PlayFolder(c *config.Connection, uris ...string) {
	cm := c.Client.BeginCommandList()
	cm.Clear()
	for _, uri := range uris {
		cm.Add(uri)
	}
	cm.Play(0)
	err := cm.End()
	config.Log(err)
}

// adds the folder to the current queue
func AddFolder(c *config.Connection, position string, uris ...string) {
	add := func(p string) {
		cm := c.Client.BeginCommandList()
		for _, uri := range uris {
			cm.Add(uri)
			c.Client.Command("add %s %s", uris, p)
		}
		err := cm.End()
		config.Log(err)
	}
	switch position {
	case "currentSong":
		add("+0")
	case "endOfQueue":
		add("")
	case "currentAlbum":
		cs, err := c.Client.CurrentSong()
		if err != nil {
			config.Log(err)
			return
		}
		pos, err := strconv.Atoi(cs["pos"])
		if err != nil {
			config.Log(err)
			return
		}
		_, endOfAlbum := findSongsAlbum(c, pos)
		add(fmt.Sprint(endOfAlbum))
	}
}

func newFolder(item mpd.Attrs) Folder {
	return Folder{
		Directory: item["directory"],
	}
}

func newFile(item mpd.Attrs) File {
	return File{
		File: item["file"],
	}
}
