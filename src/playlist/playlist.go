package playlist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

var err error

type Queue struct {
	Length   int
	Duration int
	Albums   []Album
}
type Album struct {
	Songs               []Song
	Album, Artist, Date string
}
type Song struct {
	Title, Pos, Track, Duration string
}

//returns current queue list
func GetQueue(c config.Connection) Queue {
	albumsCount := 0
	queueDuration := 0.0
	c.Connect()
	queue, err := c.Client.PlaylistInfo(-1, -1)
	c.Close()
	config.Log(err)
	for _, song := range queue {
		duration, _ := strconv.ParseFloat(song["duration"], 64)
		queueDuration += duration
	}
	length := len(queue)
	if length > 500 {
		queue = limitQueue(c, queue)
	} else if length == 0 {
		return Queue{Length: 0}
	}
	Albums := make([]Album, 0)
	Albums = append(Albums, newAlbum(queue[0]))
	Albums[0].Songs = append(Albums[0].Songs, newSong(queue[0]))
	for i := 1; i < len(queue); i++ {
		if queue[i]["Album"] != queue[i-1]["Album"] {
			albumsCount++
			Albums = append(Albums, newAlbum(queue[i]))
		}
		Albums[albumsCount].Songs = append(Albums[albumsCount].Songs, newSong(queue[i]))
	}
	return Queue{
		Albums:   Albums,
		Length:   length,
		Duration: int(queueDuration),
	}
}

//loads the playlist after the current song in queue
func AddAfterCurrent(c *config.Connection, name string) {
	err = c.Client.Command("load %s 0: +0", name).OK()
}

//removes duplicate songs based on their file address from the playlist name
//if name is empty then it deletes duplicate songs in current queue
func RemoveDuplicatesongs(c *config.Connection, name string) {
	var queue []mpd.Attrs
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
			if name == "" {
				index, _ := strconv.Atoi(queue[i]["Pos"])
				cmds.Delete(index, -1)
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

//removes songs that are removed or replaced from the server
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

//plays the id song in current Queue
func PlaySong(c *config.Connection, id int) {
	err := c.Client.Play(id)
	config.Log(err)
}

//moves the song from position in queue to newPosition
func MoveSong(c *config.Connection, position int, newPosition int) {
	err := c.Client.Move(position, -1, newPosition)
	config.Log(err)
}

//deletes the song from start to end from current Queue
func DeleteSong(c *config.Connection, start, end int) {
	err := c.Client.Delete(start, end)
	config.Log(err)
}

//saves the current queue as a new playlist
func SaveQueue(c *config.Connection, name string) {
	err := c.Client.PlaylistSave(name)
	config.Log(err)
}

//return a list of playlists
func ListStoredPlaylist(c config.Connection) (playlist []mpd.Attrs) {
	c.Connect()
	playlist, err := c.Client.ListPlaylists()
	config.Log(err)
	for _, v := range playlist {
		playlistDuration := 0.0
		songs, _ := c.Client.PlaylistContents(v["playlist"])
		for _, song := range songs {
			duration, _ := strconv.ParseFloat(song["duration"], 64)
			playlistDuration += duration
		}
		v["songsCount"] = fmt.Sprintf("%d", len(songs))
		v["duration"] = fmt.Sprintf("%f", playlistDuration)
	}
	c.Close()
	return
}

// returns list of playlist's song
func ListSongs(c *config.Connection, playlist string) (songs []mpd.Attrs) {
	contents, err := c.Client.PlaylistContents(playlist)
	config.Log(err)
	for _, song := range contents {
		m := map[string]string{
			"Title":  song["Title"],
			"Album":  song["Album"],
			"Artist": song["Artist"]}
		songs = append(songs, m)
	}
	return
}

//clears the stored playlist from server
func ClearPlaylist(c *config.Connection, playlist string) {
	err := c.Client.PlaylistClear(playlist)
	config.Log(err)
}

//deletes the stored playlist from server
func DeletePlaylist(c *config.Connection, playlist string) {
	err := c.Client.Command("rm %s", playlist).OK()
	config.Log(err)
}

//adds the playlist to the current queue
func LoadPlaylist(c *config.Connection, playlist string) {
	err := c.Client.PlaylistLoad(playlist, -1, -1)
	config.Log(err)
}

//clears the current queue and plays the playlist
func PlayPlaylist(c *config.Connection, name string) {
	cm := c.Client.BeginCommandList()
	cm.Clear()
	cm.PlaylistLoad(name, -1, -1)
	cm.Play(0)
	err := cm.End()
	config.Log(err)
}

//renames the name playlist to newName
func RenamePlaylist(c *config.Connection, name, newName string) {
	err := c.Client.PlaylistRename(name, newName)
	config.Log(err)
}

//returns content of the folder
func ListFolders(c *config.Connection, folder string) (folderAndFiles []mpd.Attrs) {
	files := make([]mpd.Attrs, 0)
	contents, err := c.Client.ListInfo(folder)
	config.Log(err)
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

//adds the folder to the current queue
func AddFolder(c *config.Connection, uris ...string) {
	cm := c.Client.BeginCommandList()
	for _, uri := range uris {
		cm.Add(uri)
	}
	err := cm.End()
	config.Log(err)
}

//adds the song to playlist
func AddSongToPlaylist(c *config.Connection, playlist, uri string) {
	err := c.Client.PlaylistAdd(playlist, uri)
	config.Log(err)
}

//search for songs in the server
func SearchServer(c *config.Connection, term ...string) (files []mpd.Attrs) {
	files, err = c.Client.Search(term...)
	config.Log(err)
	if len(files) >= 100 {
		files = files[:100]
	}
	return files
}

//returns a list of liked songs
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

//returns the list of most played songs
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

//gets the file path of the Pos in Queue
func GetSongFile(c *config.Connection, Pos int) string {
	song, err := c.Client.PlaylistInfo(Pos, -1)
	config.Log(err)
	return song[0]["file"]
}

//filter album info
func newAlbum(song mpd.Attrs) Album {
	return Album{
		Album:  song["Album"],
		Artist: song["Artist"],
		Date:   song["Date"],
	}
}

//filter song info
func newSong(song mpd.Attrs) Song {
	return Song{
		Title:    song["Title"],
		Pos:      song["Pos"],
		Track:    song["Track"],
		Duration: song["duration"],
	}
}

func limitQueue(c config.Connection, queue []mpd.Attrs) []mpd.Attrs {
	c.Connect()
	cr, err := c.Client.CurrentSong()
	c.Close()
	config.Log(err)
	cu := cr["Pos"]
	po, _ := strconv.Atoi(cu)
	if po >= 150 {
		return queue[po-150 : po+150]
	}
	return queue[:300]
}
