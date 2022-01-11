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
	if len(queue) == 0 {
		return nil
	}
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

//removes duplicate songs based on their file address from the playlist name
//if name is empty then it deletes duplicate songs in current queue
func RemoveDuplicatesongs(c config.Connection, name string) {
	var queue []mpd.Attrs
	c.Connect()
	if name == "" {
		queue, err = c.Client.PlaylistInfo(-1, -1)
		config.Log(err)
	} else {
		queue, err = c.Client.PlaylistContents(name)
		config.Log(err)
	}
	var songs = make(map[string]bool)
	cmds := c.Client.BeginCommandList()
	for i := len(queue) - 1; i >= 0; i-- {
		if _, ok := songs[queue[i]["file"]]; ok {
			index, _ := strconv.Atoi(queue[i]["Pos"])
			if name == "" {
				cmds.Delete(index, -1)
				fmt.Println(index)
			} else {
				cmds.PlaylistDelete(name, index)
			}
			continue
		}
		songs[queue[i]["file"]] = true
	}
	err := cmds.End()
	config.Log(err)
}

//plays the id song in current Queue
func PlaySong(c config.Connection, id int) {
	c.Connect()
	err := c.Client.Play(id)
	config.Log(err)
	c.Close()
}

//moves the song from position in queue to newPosition
func MoveSong(c config.Connection, position int, newPosition int) {
	c.Connect()
	err := c.Client.Move(position, -1, newPosition)
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

//return a list of playlists
func ListStoredPlaylist(c config.Connection) (playlist []mpd.Attrs) {
	c.Connect()
	playlist, err := c.Client.ListPlaylists()
	config.Log(err)
	for _, v := range playlist {
		songCount := 0
		playlistDuration := 0.0
		songs, _ := c.Client.PlaylistContents(v["playlist"])
		for _, song := range songs {
			duration, _ := strconv.ParseFloat(song["duration"], 64)
			playlistDuration += duration
			songCount++
		}
		v["songsCount"] = fmt.Sprintf("%d", songCount)
		v["duration"] = fmt.Sprintf("%f", playlistDuration)
	}
	c.Close()
	return
}

// returns list of playlist's song
func ListSongs(c config.Connection, playlist string) (songs []mpd.Attrs) {
	c.Connect()
	contents, err := c.Client.PlaylistContents(playlist)
	config.Log(err)
	for _, song := range contents {
		m := map[string]string{
			"Title":  song["Title"],
			"Album":  song["Album"],
			"Artist": song["Artist"]}
		songs = append(songs, m)
	}
	c.Close()
	return
}

//clears the stored playlist from server
func ClearPlaylist(c config.Connection, playlist string) {
	c.Connect()
	err := c.Client.PlaylistClear(playlist)
	config.Log(err)
	c.Close()
}

//deletes the stored playlist from server
func DeletePlaylist(c config.Connection, playlist string) {
	c.Connect()
	err := c.Client.Command("rm %s", playlist).OK()
	config.Log(err)
	c.Close()
}

//adds the playlist to the current queue
func LoadPlaylist(c config.Connection, playlist string) {
	c.Connect()
	err := c.Client.PlaylistLoad(playlist, -1, -1)
	config.Log(err)
	c.Close()
}

//clears the current queue and plays the playlist
func PlayPlaylist(c config.Connection, name string) {
	c.Connect()
	cm := c.Client.BeginCommandList()
	cm.Clear()
	cm.PlaylistLoad(name, -1, -1)
	cm.Play(0)
	err := cm.End()
	config.Log(err)
	c.Close()
}

//renames the name playlist to newName
func RenamePlaylist(c config.Connection, name, newName string) {
	c.Connect()
	err := c.Client.PlaylistRename(name, newName)
	config.Log(err)
	c.Close()
}

//returns content of the folder
func ListFolders(c config.Connection, folder string) (folderAndFiles []mpd.Attrs) {
	files := make([]mpd.Attrs, 0)
	c.Connect()
	contents, err := c.Client.ListInfo(folder)
	config.Log(err)
	c.Close()
	for _, item := range contents {
		//ignore playlists
		if _, ok := item["playlist"]; ok {
			continue
		}
		//append files to the end of list
		if _, ok := item["file"]; ok {
			files = append(files, item)
			continue
		}
		folderAndFiles = append(folderAndFiles, item)
	}
	folderAndFiles = append(folderAndFiles, files...)
	return
}

//clears the current queue and plays the folder
func PlayFolder(c config.Connection, uris ...string) {
	c.Connect()
	cm := c.Client.BeginCommandList()
	cm.Clear()
	for _, uri := range uris {
		cm.Add(uri)
	}
	cm.Play(0)
	err := cm.End()
	config.Log(err)
	c.Close()
}

//adds the folder to the current queue
func AddFolder(c config.Connection, uris ...string) {
	c.Connect()
	cm := c.Client.BeginCommandList()
	for _, uri := range uris {
		cm.Add(uri)
	}
	err := cm.End()
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

//search for songs in the server
func SearchServer(c config.Connection, term ...string) (files []mpd.Attrs) {
	c.Connect()
	files, err = c.Client.Search(term...)
	config.Log(err)
	c.Close()
	if len(files) >= 100 {
		files = files[:100]
	}
	return files
}

//returns a list of liked songs
func GetLikedSongs(c config.Connection) (likedSongs []string) {
	c.Connect()
	uris, stickers, e := c.Client.StickerFind("", "liked")
	config.Log(e)
	c.Close()
	likedSongs = make([]string, 0)
	for i := 0; i < len(stickers); i++ {
		if stickers[i].Value == "true" {
			likedSongs = append(likedSongs, uris[i])
		}
	}
	return likedSongs
}

//returns the list of most played songs
func GetMostPlayed(c config.Connection) (mostPlayed []string) {
	c.Connect()
	uris, stickers, err := c.Client.StickerFind("", "playedcount")
	c.Close()
	config.Log(err)
	count := len(uris)
	unOrdered := make([]map[string]string, 0)
	mostPlayed = make([]string, 0)
	for i := 0; i < count; i++ {
		m := map[string]string{"uri": uris[i], "count": stickers[i].Value}
		unOrdered = append(unOrdered, m)
	}
	//sort songs based on played time
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
