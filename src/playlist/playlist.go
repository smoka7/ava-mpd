package playlist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

var err error

//returns current queue list
func GetQueue(c config.Connection) [][]mpd.Attrs {
	c.Connect()
	newIndex := 0
	queue, err := c.Client.PlaylistInfo(-1, -1)
	config.Log(err)
	c.Close()
	filteredQueue := make([][]mpd.Attrs, 1)
	filteredQueue[0] = append(filteredQueue[0], newAlbum(queue[0]))
	for i := 1; i < len(queue); i++ {
		if queue[i]["Album"] == queue[i-1]["Album"] {
			filteredQueue[newIndex] = append(filteredQueue[newIndex], newSong(queue[i]))
			continue
		}
		newIndex++
		//increase the number of albums
		filteredQueue = append(filteredQueue, []mpd.Attrs{})
		filteredQueue[newIndex] = append(filteredQueue[newIndex], newAlbum(queue[i]))
	}
	return filteredQueue
}

//plays the id song in current Queue
func PlaySong(c config.Connection, id int) {
	c.Connect()
	err := c.Client.Play(id)
	config.Log(err)
	c.Close()
}

//deletes the song from start to end from current Queue
func DeleteSong(c config.Connection, start, end int) {
	c.Connect()
	err := c.Client.Delete(start, end)
	config.Log(err)
	c.Close()
}

//saves the current queue as a new playlist
func SaveQueue(c config.Connection, name string) {
	c.Connect()
	err := c.Client.PlaylistSave(name)
	config.Log(err)
	c.Close()
}

//adds the song to playlist
func AddSongToPlaylist(c config.Connection, playlist, uri string) {
	c.Connect()
	err := c.Client.PlaylistAdd(playlist, uri)
	config.Log(err)
	c.Close()
}

//filter album info
func newAlbum(song mpd.Attrs) mpd.Attrs {
	return map[string]string{
		"Title":    song["Title"],
		"Album":    song["Album"],
		"Pos":      song["Pos"],
		"Track":    song["Track"],
		"Artist":   song["Artist"],
		"duration": song["duration"],
		"Date":     song["Date"],
		"file":     song["file"],
	}
}

//filter song info
func newSong(song mpd.Attrs) mpd.Attrs {
	return map[string]string{
		"Title":    song["Title"],
		"Pos":      song["Pos"],
		"Track":    song["Track"],
		"duration": song["duration"],
		"file":     song["file"],
	}
}
