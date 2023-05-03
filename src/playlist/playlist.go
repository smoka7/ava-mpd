package playlist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

type Request struct {
	Command string `json:"command"`
	Data    Data   `json:"data"`
}

type Data struct {
	Playlist    string
	NewPlaylist string
	Pos         string
}

// loads the playlist after the current song in queue.
func (a action) addAfterPos(name, pos string) {
	// 0: means load whole playlist
	err := a.Command("load %s 0: %s", name, pos).OK()
	config.Log(err)
}

// removes duplicate songs based on their file address from the playlist name
// if name is empty then it deletes duplicate songs in current queue.
func (a action) RemoveDuplicatesongs(name string) {
	removeFromQueue := name == ""
	getQueue := func() []mpd.Attrs {
		if removeFromQueue {
			queue, err := a.PlaylistInfo(-1, -1)
			config.Log(err)

			return queue
		}

		queue, err := a.PlaylistContents(name)
		config.Log(err)
		return queue
	}

	queue := getQueue()
	cmds := a.BeginCommandList()

	// cmd used to delete duplicates
	var removeDuplicates func(int)
	if removeFromQueue {
		removeDuplicates = func(index int) {
			id, _ := strconv.Atoi(queue[index]["Id"])
			cmds.DeleteID(id)
		}
	} else {
		removeDuplicates = func(i int) {
			cmds.PlaylistDelete(name, i)
		}
	}

	// find the duplicate songs
	songs := make(map[string]struct{})
	for i := len(queue) - 1; i >= 0; i-- {
		if _, duplicate := songs[queue[i]["file"]]; duplicate {
			removeDuplicates(i)

			continue
		}
		songs[queue[i]["file"]] = struct{}{}
	}

	err := cmds.End()
	config.Log(err)
}

// removes songs that are removed or replaced from the server.
func (a action) RemoveInvalidsongs(name string) {
	queue, err := a.PlaylistContents(name)
	config.Log(err)
	cmds := a.BeginCommandList()
	for i := len(queue) - 1; i >= 0; i-- {
		if _, ok := queue[i]["Title"]; !ok {
			cmds.PlaylistDelete(name, i)
		}
	}
	err = cmds.End()
	config.Log(err)
}

// return a list of playlists.
func (a action) ListStoredPlaylist() (playlists Playlists) {
	a.Connect()
	defer a.Close()
	list, err := a.ListPlaylists()
	config.Log(err)
	for _, v := range list {
		songs, _ := a.PlaylistContents(v["playlist"])
		playlists = append(playlists, Playlist{
			Name:       v["playlist"],
			Duration:   getDurationSum(songs),
			SongsCount: uint(len(songs)),
		})
	}

	return
}

// returns list of playlist's song.
func (a action) ListSongsIn(playlist string) (songs Songs) {
	contents, err := a.PlaylistContents(playlist)
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

// clears the stored playlist from server.
func (a action) ClearPlaylist(playlist string) {
	err := a.PlaylistClear(playlist)
	config.Log(err)
}

// deletes the stored playlist from server.
func (a action) DeletePlaylist(playlist string) {
	err := a.Command("rm %s", playlist).OK()
	config.Log(err)
}

// adds the playlist to the current queue.
func (a action) LoadPlaylist(d Data) {
	pos := a.getCurrentSongPos()
	if pos < 0 {
		err := a.PlaylistLoad(d.Playlist, -1, -1)
		config.Log(err)
		return
	}

	switch d.Pos {
	case "currentSong":
		a.addAfterPos(d.Playlist, "+0")

	case "endOfQueue":
		err := a.PlaylistLoad(d.Playlist, -1, -1)
		config.Log(err)

	case "currentAlbum":
		_, endOfAlbum := a.findSongsAlbum(pos)
		a.addAfterPos(d.Playlist, fmt.Sprint(endOfAlbum))
	}
}

// clears the current queue and plays the playlist.
func (a action) PlayPlaylist(d Data) {
	cm := a.BeginCommandList()
	cm.Clear()
	cm.PlaylistLoad(d.Playlist, -1, -1)
	cm.Play(0)
	err := cm.End()
	config.Log(err)
}

// renames the name playlist to newName.
func (a action) RenamePlaylist(d Data) {
	err := a.PlaylistRename(d.Playlist, d.NewPlaylist)
	config.Log(err)
}

// adds the song to playlist.
func (a action) AddSongToPlaylist(playlist, file string) {
	err := a.PlaylistAdd(playlist, file)
	config.Log(err)
}

// returns a list of liked songs.
func (a action) GetLikedSongs() []string {
	uris := getLikedURIS(a.Connection)
	likedSongs := make([]string, len(uris))
	i := 0
	for uri := range uris {
		likedSongs[i] = uri
		i++
	}

	sort.Slice(likedSongs, func(i, j int) bool {
		return likedSongs[i] < likedSongs[j]
	})

	return likedSongs
}

type song struct {
	uri         string
	playedCount string
}

// returns the list of most played songs.
func (a action) GetMostPlayed() []string {
	uris, stickers, err := a.StickerFind("", "playedcount")
	config.Log(err)
	count := len(uris)

	unOrdered := make([]song, count)
	for i := 0; i < count; i++ {
		unOrdered[i] = song{uri: uris[i], playedCount: stickers[i].Value}
	}

	// sort songs based on played time
	sort.Slice(unOrdered, func(i, j int) bool {
		a, _ := strconv.Atoi(unOrdered[i].playedCount)
		b, _ := strconv.Atoi(unOrdered[j].playedCount)
		return a > b
	})

	if count > 100 {
		count = 100
	}
	mostPlayed := make([]string, count)

	for i := 0; i < count; i++ {
		mostPlayed[i] = unOrdered[i].uri
	}
	return mostPlayed
}

func getLikedURIS(c config.Connection) map[string]struct{} {
	c.Connect()
	attrs, err := c.Command("sticker find song %s %s = %s", "", "liked", "true").AttrsList("file")
	c.Close()
	if err != nil {
		return nil
	}
	config.Log(err)
	uris := make(map[string]struct{}, len(attrs))
	for _, song := range attrs {
		uris[song["file"]] = struct{}{}
	}
	return uris
}

// returns the duration of a list of songs.
func getDurationSum(songs []mpd.Attrs) uint {
	sum := 0.0
	for _, song := range songs {
		duration, _ := strconv.ParseFloat(song["duration"], 64)
		sum += duration
	}

	return uint(sum)
}
