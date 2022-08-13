package playlist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

type PlaylistRequest struct {
	Command string       `json:"command"`
	Data    PlaylistData `json:"data"`
}

type PlaylistData struct {
	Playlist    string
	NewPlaylist string
	Pos         string
}

// loads the playlist after the current song in queue
func (a action) addAfterPos(name, pos string) {
	// 0: means load whole playlist
	err = a.c.Client.Command("load %s 0: %s", name, pos).OK()
}

// removes duplicate songs based on their file address from the playlist name
// if name is empty then it deletes duplicate songs in current queue
func (a action) RemoveDuplicatesongs(name string) {
	getQueue := func() (queue []mpd.Attrs) {
		if name == "" {
			queue, err = a.c.Client.PlaylistInfo(-1, -1)
			config.Log(err)
			return
		}
		queue, err = a.c.Client.PlaylistContents(name)
		config.Log(err)
		return
	}

	queue := getQueue()
	songs := make(map[string]bool)
	cmds := a.c.Client.BeginCommandList()
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
func (a action) RemoveInvalidsongs(name string) {
	queue, err := a.c.Client.PlaylistContents(name)
	config.Log(err)
	cmds := a.c.Client.BeginCommandList()
	for i := len(queue) - 1; i >= 0; i-- {
		if _, ok := queue[i]["Title"]; !ok {
			cmds.PlaylistDelete(name, i)
		}
	}
	err = cmds.End()
	config.Log(err)
}

// return a list of playlists
func (a action) ListStoredPlaylist() (playlists Playlists) {
	a.c.Connect()
	defer a.c.Close()
	list, err := a.c.Client.ListPlaylists()
	config.Log(err)
	for _, v := range list {
		songs, _ := a.c.Client.PlaylistContents(v["playlist"])
		playlists = append(playlists, Playlist{
			Name:       v["playlist"],
			Duration:   getDurationSum(songs),
			SongsCount: uint(len(songs)),
		})
	}
	return
}

// returns list of playlist's song
func (a action) ListSongsIn(playlist string) (songs Songs) {
	contents, err := a.c.Client.PlaylistContents(playlist)
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
func (a action) ClearPlaylist(playlist string) {
	err := a.c.Client.PlaylistClear(playlist)
	config.Log(err)
}

// deletes the stored playlist from server
func (a action) DeletePlaylist(playlist string) {
	err := a.c.Client.Command("rm %s", playlist).OK()
	config.Log(err)
}

// adds the playlist to the current queue
func (a action) LoadPlaylist(d PlaylistData) {
	switch d.Pos {
	case "currentSong":
		a.addAfterPos(d.Playlist, "+0")
	case "endOfQueue":
		err := a.c.Client.PlaylistLoad(d.Playlist, -1, -1)
		config.Log(err)
	case "currentAlbum":
		pos := a.getCurrentSongPos()
		_, endOfAlbum := a.findSongsAlbum(pos)
		a.addAfterPos(d.Playlist, fmt.Sprint(endOfAlbum))
	}
}

// clears the current queue and plays the playlist
func (a action) PlayPlaylist(d PlaylistData) {
	cm := a.c.Client.BeginCommandList()
	cm.Clear()
	cm.PlaylistLoad(d.Playlist, -1, -1)
	cm.Play(0)
	err := cm.End()
	config.Log(err)
}

// renames the name playlist to newName
func (a action) RenamePlaylist(d PlaylistData) {
	err := a.c.Client.PlaylistRename(d.Playlist, d.NewPlaylist)
	config.Log(err)
}

// adds the song to playlist
func (a action) AddSongToPlaylist(playlist, file string) {
	err := a.c.Client.PlaylistAdd(playlist, file)
	config.Log(err)
}

// returns a list of liked songs
func (a action) GetLikedSongs() []string {
	uris, stickers, e := a.c.Client.StickerFind("", "liked")
	config.Log(e)
	likedSongs := make([]string, 0)
	for i := 0; i < len(stickers); i++ {
		if stickers[i].Value == "true" {
			likedSongs = append(likedSongs, uris[i])
		}
	}
	return likedSongs
}

// returns the list of most played songs
func (a action) GetMostPlayed() []string {
	uris, stickers, err := a.c.Client.StickerFind("", "playedcount")
	config.Log(err)
	count := len(uris)
	unOrdered := make([]map[string]string, 0)
	mostPlayed := make([]string, 0)
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
	return mostPlayed
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
