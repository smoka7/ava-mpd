package playlist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

// returns current queue list
func GetQueue(c config.Connection) (q Queue) {
	c.Connect()
	queue, err := c.Client.PlaylistInfo(-1, -1)
	c.Close()
	config.Log(err)
	q.Duration = getDurationSum(queue)
	q.Length = len(queue)
	if q.Length == 0 {
		return
	} else if q.Length > 500 {
		queue = limitQueue(c, queue)
	}
	q.newAlbum(queue[0])
	for i := 1; i < len(queue); i++ {
		if queue[i]["Album"] != queue[i-1]["Album"] {
			q.newAlbum(queue[i])
			continue
		}
		q.Albums[len(q.Albums)-1].newSong(queue[i])
	}
	return
}

// loads the playlist after the current song in queue
func AddAfterCurrent(c *config.Connection, name string) {
	err = c.Client.Command("load %s 0: +0", name).OK()
}

// removes duplicate songs based on their file address from the playlist name
// if name is empty then it deletes duplicate songs in current queue
func RemoveDuplicatesongs(c *config.Connection, name string) {
	var queue []mpd.Attrs
	if name == "" {
		queue, err = c.Client.PlaylistInfo(-1, -1)
		config.Log(err)
	} else {
		queue, err = c.Client.PlaylistContents(name)
		config.Log(err)
	}
	songs := make(map[string]bool)
	cmds := c.Client.BeginCommandList()
	for i := len(queue) - 1; i >= 0; i-- {
		if _, ok := songs[queue[i]["file"]]; ok {
			if name == "" {
				id, _ := strconv.Atoi(queue[i]["Id"])
				cmds.DeleteID(id)
			} else {
				cmds.PlaylistDelete(name, i)
			}
			continue
		}
		songs[queue[i]["file"]] = true
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

// plays the id song in current Queue
func PlaySong(c *config.Connection, id int) {
	err := c.Client.PlayID(id)
	config.Log(err)
}

// moves the song from position in queue to newPosition
func MoveSong(c *config.Connection, position int, newPosition int) {
	err := c.Client.MoveID(position, newPosition)
	config.Log(err)
}

// deletes the song from start to end from current Queue

func DeleteSong(c *config.Connection, id int) {
	err := c.Client.DeleteID(id)
	// err := c.Client.Delete(start, end)
	config.Log(err)
}

// saves the current queue as a new playlist
func SaveQueue(c *config.Connection, name string) {
	err := c.Client.PlaylistSave(name)
	config.Log(err)
}

// return a list of playlists
func ListStoredPlaylist(c config.Connection) (playlist []mpd.Attrs) {
	c.Connect()
	playlist, err := c.Client.ListPlaylists()
	config.Log(err)
	for _, v := range playlist {
		songs, _ := c.Client.PlaylistContents(v["playlist"])
		v["songsCount"] = fmt.Sprintf("%d", len(songs))
		v["duration"] = fmt.Sprintf("%f", getDurationSum(songs))
	}
	c.Close()
	return
}

//  returns list of playlist's song
func ListSongs(c *config.Connection, playlist string) (songs []mpd.Attrs) {
	contents, err := c.Client.PlaylistContents(playlist)
	config.Log(err)
	for _, song := range contents {
		m := map[string]string{
			"Title":  song["Title"],
			"Album":  song["Album"],
			"Artist": song["Artist"],
		}
		songs = append(songs, m)
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
func LoadPlaylist(c *config.Connection, playlist string) {
	err := c.Client.PlaylistLoad(playlist, -1, -1)
	config.Log(err)
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

// gets the file path of the Pos in Queue
func GetSongFile(c *config.Connection, id int) string {
	song, err := c.Client.Command("playlistid %d", id).Attrs()
	config.Log(err)
	return song["file"]
}

// shuffles an album given the position of the song in queue
func ShuffleAlbum(c *config.Connection, pos int) {
	queue, err := c.Client.PlaylistInfo(-1, -1)
	if err != nil {
		config.Log(err)
		return
	}

	lastSongIndex := pos
	for ; lastSongIndex < len(queue); lastSongIndex++ {
		if queue[lastSongIndex]["Album"] != queue[pos]["Album"] {
			break
		}
	}

	firstSongIndex := pos
	for ; firstSongIndex >= 0; firstSongIndex-- {
		if queue[firstSongIndex]["Album"] != queue[pos]["Album"] {
			break
		}
	}

	err = c.Client.Shuffle(firstSongIndex+1, lastSongIndex)
	config.Log(err)
}

// filter album info
func (q *Queue) newAlbum(song mpd.Attrs) {
	album := Album{
		Album:  song["Album"],
		Artist: song["Artist"],
		Date:   song["Date"],
	}
	album.newSong(song)
	q.Albums = append(q.Albums, album)
}

// filter song info
func (a *Album) newSong(song mpd.Attrs) {
	a.Songs = append(a.Songs, Song{
		Title:    song["Title"],
		Pos:      song["Pos"],
		Id:       song["Id"],
		Track:    song["Track"],
		Duration: song["duration"],
	})
}

func limitQueue(c config.Connection, queue []mpd.Attrs) []mpd.Attrs {
	c.Connect()
	currentSong, err := c.Client.CurrentSong()
	c.Close()
	config.Log(err)
	queuePosition := currentSong["Pos"]
	if p, _ := strconv.Atoi(queuePosition); p >= 150 {
		return queue[p-150 : p+150]
	}
	return queue[:300]
}

// returns the duration of a list of songs
func getDurationSum(songs []mpd.Attrs) (sum float64) {
	for _, song := range songs {
		duration, _ := strconv.ParseFloat(song["duration"], 64)
		sum += duration
	}
	return
}
