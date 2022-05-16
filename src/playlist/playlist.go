package playlist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

// loads the playlist after the current song in queue
func addAfterPos(c *config.Connection, name, pos string) {
	// 0: means load whole playlist
	err = c.Client.Command("load %s 0: %s", name, pos).OK()
}

// removes duplicate songs based on their file address from the playlist name
// if name is empty then it deletes duplicate songs in current queue
func RemoveDuplicatesongs(c *config.Connection, name string) {
	getQueue := func() (queue []mpd.Attrs) {
		if name == "" {
			queue, err = c.Client.PlaylistInfo(-1, -1)
			config.Log(err)
			return
		}
		queue, err = c.Client.PlaylistContents(name)
		config.Log(err)
		return
	}

	queue := getQueue()
	songs := make(map[string]bool)
	cmds := c.Client.BeginCommandList()
	for i := len(queue) - 1; i >= 0; i-- {

		if _, duplicate := songs[queue[i]["file"]]; !duplicate {
			songs[queue[i]["file"]] = true
			continue
		}
		// deleting from Queue
		if name == "" {
			id, _ := strconv.Atoi(queue[i]["Id"])
			cmds.DeleteID(id)
			continue
		}
		// deleting from StoredPlaylist
		cmds.PlaylistDelete(name, i)
	}
	err := cmds.End()
	config.Log(err)
}

// removes songs that are removed or replaced from the server
func RemoveInvalidsongs(c *config.Connection, name string) {
	queue, err := c.Client.PlaylistContents(name)
	config.Log(err)
	cmds := c.Client.BeginCommandList()
	for i := len(queue) - 1; i >= 0; i-- {
		if _, ok := queue[i]["Title"]; !ok {
			cmds.PlaylistDelete(name, i)
		}
	}
	err = cmds.End()
	config.Log(err)
}

// return a list of playlists
func ListStoredPlaylist(c config.Connection) (playlists Playlists) {
	c.Connect()
	list, err := c.Client.ListPlaylists()
	config.Log(err)
	for _, v := range list {
		songs, _ := c.Client.PlaylistContents(v["playlist"])
		playlists = append(playlists, Playlist{
			Name:       v["playlist"],
			Duration:   getDurationSum(songs),
			SongsCount: uint(len(songs)),
		})
	}
	c.Close()
	return
}

//  returns list of playlist's song
func ListSongsIn(c *config.Connection, playlist string) (songs Songs) {
	contents, err := c.Client.PlaylistContents(playlist)
	config.Log(err)
	for _, song := range contents {
		songs = append(songs, Song{
			Title:  song["Title"],
			Album:  song["Album"],
			Artist: song["Artist"],
		})
	}
	return
}

// clears the stored playlist from server
func ClearPlaylist(c *config.Connection, playlist string) {
	err := c.Client.PlaylistClear(playlist)
	config.Log(err)
}

// deletes the stored playlist from server
func DeletePlaylist(c *config.Connection, playlist string) {
	err := c.Client.Command("rm %s", playlist).OK()
	config.Log(err)
}

// adds the playlist to the current queue
func LoadPlaylist(c *config.Connection, playlist, position string) {
	switch position {
	case "currentSong":
		addAfterPos(c, playlist, "+0")
	case "endOfQueue":
		err := c.Client.PlaylistLoad(playlist, -1, -1)
		config.Log(err)
	case "currentAlbum":
		pos := getCurrentSongPos(c)
		_, endOfAlbum := findSongsAlbum(c, pos)
		addAfterPos(c, playlist, fmt.Sprint(endOfAlbum))
	}
}

// clears the current queue and plays the playlist
func PlayPlaylist(c *config.Connection, name string) {
	cm := c.Client.BeginCommandList()
	cm.Clear()
	cm.PlaylistLoad(name, -1, -1)
	cm.Play(0)
	err := cm.End()
	config.Log(err)
}

// renames the name playlist to newName
func RenamePlaylist(c *config.Connection, name, newName string) {
	err := c.Client.PlaylistRename(name, newName)
	config.Log(err)
}

// adds the song to playlist
func AddSongToPlaylist(c *config.Connection, playlist, uri string) {
	err := c.Client.PlaylistAdd(playlist, uri)
	config.Log(err)
}

// returns a list of liked songs
func GetLikedSongs(c *config.Connection) (likedSongs []string) {
	uris, stickers, e := c.Client.StickerFind("", "liked")
	config.Log(e)
	likedSongs = make([]string, 0)
	for i := 0; i < len(stickers); i++ {
		if stickers[i].Value == "true" {
			likedSongs = append(likedSongs, uris[i])
		}
	}
	return likedSongs
}

// returns the list of most played songs
func GetMostPlayed(c *config.Connection) (mostPlayed []string) {
	uris, stickers, err := c.Client.StickerFind("", "playedcount")
	config.Log(err)
	count := len(uris)
	unOrdered := make([]map[string]string, 0)
	mostPlayed = make([]string, 0)
	for i := 0; i < count; i++ {
		m := map[string]string{"uri": uris[i], "count": stickers[i].Value}
		unOrdered = append(unOrdered, m)
	}
	// sort songs based on played time
	sort.Slice(unOrdered, func(i, j int) bool {
		a, err := strconv.Atoi(unOrdered[i]["count"])
		config.Log(err)
		b, err := strconv.Atoi(unOrdered[j]["count"])
		config.Log(err)
		return a > b
	})
	if count > 100 {
		count = 100
	}
	for i := 0; i < count; i++ {
		mostPlayed = append(mostPlayed, unOrdered[i]["uri"])
	}
	return
}

// returns the duration of a list of songs
func getDurationSum(songs []mpd.Attrs) uint {
	sum := 0.0
	for _, song := range songs {
		duration, _ := strconv.ParseFloat(song["duration"], 64)
		sum += duration
	}
	return uint(sum)
}
