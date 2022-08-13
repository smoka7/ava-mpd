package playlist

import (
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

type QueueRequest struct {
	Command string    `json:"command"`
	Data    QueueData `json:"data"`
}

type QueueData struct {
	ID, Start, Finish int
	Playlist          string
}

type action struct {
	c config.Connection
}

const ClientQueueLimit = 200

func NewAction(c config.Connection) action {
	if c.Client != nil {
		return action{c}
	}
	c.Connect()
	defer c.Connect()
	return action{c}
}

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
func (a action) PlaySong(d QueueData) {
	err := a.c.Client.PlayID(d.ID)
	config.Log(err)
}

// deletes the song from start to end from current Queue
func (a action) DeleteSong(d QueueData) {
	err := a.c.Client.DeleteID(d.ID)
	config.Log(err)
}

// moves the song from position in queue to newPosition
func (a action) MoveSong(d QueueData) {
	err := a.c.Client.MoveID(d.Start, d.Finish)
	config.Log(err)
}

// shuffles an album given the position of the song in queue
func (a action) ShuffleAlbum(d QueueData) {
	firstSongIndex, lastSongIndex := a.findSongsAlbum(d.Start)

	err = a.c.Client.Shuffle(firstSongIndex+1, lastSongIndex)
	config.Log(err)
}

// saves the current queue as a new playlist
func (a action) SaveQueue(d QueueData) {
	err := a.c.Client.PlaylistSave(d.Playlist)
	config.Log(err)
}

// gets the file path of song with id in Queue
func (a action) GetSongFile(id int) (string, error) {
	song, err := a.c.Client.Command("playlistid %d", id).Attrs()
	if err != nil {
		return "", err
	}
	return song["file"], nil
}

// returns the boundaries of song in Queue
func (a action) findSongsAlbum(pos int) (int, int) {
	queue, err := a.c.Client.PlaylistInfo(-1, -1)
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
func (a action) getCurrentSongPos() int {
	cs, err := a.c.Client.CurrentSong()
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
