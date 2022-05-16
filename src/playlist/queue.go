package playlist

import (
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

// plays the id song in current Queue
func PlaySong(c *config.Connection, id int) {
	err := c.Client.PlayID(id)
	config.Log(err)
}

// deletes the song from start to end from current Queue
func DeleteSong(c *config.Connection, id int) {
	err := c.Client.DeleteID(id)
	config.Log(err)
}

// moves the song from position in queue to newPosition
func MoveSong(c *config.Connection, position int, newPosition int) {
	err := c.Client.MoveID(position, newPosition)
	config.Log(err)
}

// shuffles an album given the position of the song in queue
func ShuffleAlbum(c *config.Connection, pos int) {
	firstSongIndex, lastSongIndex := findSongsAlbum(c, pos)

	err = c.Client.Shuffle(firstSongIndex+1, lastSongIndex)
	config.Log(err)
}

// saves the current queue as a new playlist
func SaveQueue(c *config.Connection, name string) {
	err := c.Client.PlaylistSave(name)
	config.Log(err)
}

// gets the file path of song with id in Queue
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
	parsedPage, err := strconv.Atoi(page)
	if parsedPage == -1 {
		q.CurrentPage = q.CurrentSongPage
	} else {
		q.CurrentPage = uint(parsedPage)
	}

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
