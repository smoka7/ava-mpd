package playlist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

const ClientQueueLimit = 200

// returns current queue list
func GetQueue(c config.Connection, page string) (q Queue) {
	c.Connect()
	queue, err := c.Client.PlaylistInfo(-1, -1)
	c.Close()
	config.Log(err)
	q.Length = uint(len(queue))
	if q.Length == 0 {
		return
	}
	q.Duration = getDurationSum(queue)
	queue = q.getQueueInPage(c, queue, page)
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
	config.Log(err)
}

// saves the current queue as a new playlist
func SaveQueue(c *config.Connection, name string) {
	err := c.Client.PlaylistSave(name)
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

// gets the file path of the Pos in Queue
func GetSongFile(c *config.Connection, id int) (string, error) {
	song, err := c.Client.Command("playlistid %d", id).Attrs()
	if err != nil {
		return "", err
	}
	return song["file"], nil
}

// returns the boundaries of song in Queue
func findSongsAlbum(c *config.Connection, pos int) (int, int) {
	queue, err := c.Client.PlaylistInfo(-1, -1)
	if err != nil {
		config.Log(err)
		return 0, 0
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
	return firstSongIndex, lastSongIndex
}

// returns the current song position in queue
func getCurrentSongPos(c *config.Connection) int {
	cs, err := c.Client.CurrentSong()
	// check for empty queue
	if err != nil || cs == nil {
		config.Log(err)
		return 0
	}

	pos, err := strconv.Atoi(cs["Pos"])
	if err != nil {
		config.Log(err)
		return 0
	}
	return pos
}

// shuffles an album given the position of the song in queue
func ShuffleAlbum(c *config.Connection, pos int) {
	firstSongIndex, lastSongIndex := findSongsAlbum(c, pos)

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

func (q *Queue) getQueueInPage(c config.Connection, queue []mpd.Attrs, page string) (currentPage []mpd.Attrs) {
	c.Connect()
	currentSong, err := c.Client.CurrentSong()
	c.Close()
	config.Log(err)

	downLimit, upLimt := 0, ClientQueueLimit
	csPos, _ := strconv.Atoi(currentSong["Pos"])
	q.CurrentSongPage = uint(csPos/ClientQueueLimit + 1)
	q.LastPage = uint(len(queue)/ClientQueueLimit + 1)

	// returns the part of queue that page number requested
	parsedPage, err := strconv.ParseUint(page, 10, 32)
	q.CurrentPage = uint(parsedPage)
	if q.CurrentPage > q.LastPage {
		q.CurrentPage = q.LastPage
	} else if q.CurrentPage <= 0 {
		q.CurrentPage = 1
	}
	if err == nil || q.CurrentPage != 0 {
		downLimit = int((q.CurrentPage - 1) * ClientQueueLimit)
		upLimt = downLimit + ClientQueueLimit
	}

	if upLimt > len(queue) {
		upLimt = len(queue)
	}
	return queue[downLimit:upLimt]
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
