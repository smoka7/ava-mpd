package playlist

import (
	"strconv"
	"strings"

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
	config.Connection
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
func GetQueue(c config.Connection, pageNumber string) (q Queue) {
	c.Connect()
	queue, err := c.PlaylistInfo(-1, -1)
	c.Close()
	config.Log(err)
	q.Length = uint(len(queue))
	if q.Length == 0 {
		return
	}
	q.Duration = getDurationSum(queue)
	d, u := q.getQueueInPage(c, pageNumber)
	page := queue[d:u]
	q.newAlbum(page[0])
	liked := getLikedURIS(c)

	for i := 1; i < len(page); i++ {

		// find liked songs
		if _, ok := liked[page[i]["file"]]; ok {
			page[i]["liked"] = ""
		}

		// group album's songs together
		if page[i]["Album"] != page[i-1]["Album"] {
			q.newAlbum(page[i])
			continue
		}
		q.Albums[len(q.Albums)-1].newSong(page[i])
	}
	return
}

// plays the id song in current Queue
func (a action) PlaySong(d QueueData) {
	err := a.PlayID(d.ID)
	config.Log(err)
}

// deletes the song from start to end from current Queue
func (a action) DeleteSong(d QueueData) {
	err := a.DeleteID(d.ID)
	config.Log(err)
}

// moves the song from position in queue to newPosition
func (a action) MoveSong(d QueueData) {
	err := a.MoveID(d.Start, d.Finish)
	config.Log(err)
}

// shuffles an album given the position of the song in queue
func (a action) ShuffleAlbum(d QueueData) {
	firstSongIndex, lastSongIndex := a.findSongsAlbum(d.Start)

	err = a.Shuffle(firstSongIndex+1, lastSongIndex)
	config.Log(err)
}

// saves the current queue as a new playlist
func (a action) SaveQueue(d QueueData) {
	err := a.PlaylistSave(d.Playlist)
	config.Log(err)
}

// gets the file path of song with id in Queue
func (a action) GetSongFile(id int) (string, error) {
	song, err := a.Command("playlistid %d", id).Attrs()
	if err != nil {
		return "", err
	}
	return song["file"], nil
}

// returns the boundaries of song in Queue
func (a action) findSongsAlbum(pos int) (int, int) {
	queue, err := a.PlaylistInfo(-1, -1)
	if err != nil || len(queue) == 0 {
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
	cs, err := a.CurrentSong()
	// check for empty queue
	if err != nil {
		config.Log(err)
		return -1
	}

	if len(cs) == 0 {
		return -1
	}

	pos, err := strconv.Atoi(cs["Pos"])
	if err != nil {
		config.Log(err)
		return -1
	}
	return pos
}

// filter album info
func (q *Queue) newAlbum(song mpd.Attrs) {
	artist, ok := song["AlbumArtist"]
	if !ok {
		artist = song["Artist"]
	}

	album := Album{
		Album:  song["Album"],
		Artist: artist,
		Date:   song["Date"],
	}
	album.newSong(song)
	q.Albums = append(q.Albums, album)
}

// filter song info
func (a *Album) newSong(song mpd.Attrs) {
	if song["Title"] == "" {
		split := strings.Split(song["file"], "/")
		song["Title"] = split[len(split)-1]
	}
	_, liked := song["liked"]
	a.Songs = append(a.Songs, Song{
		Title:    song["Title"],
		Pos:      song["Pos"],
		Id:       song["Id"],
		Track:    song["Track"],
		Duration: song["duration"],
		Liked:    liked,
	})
}

func (q *Queue) getQueueInPage(c config.Connection, page string) (downLimit, upLimt int) {
	c.Connect()
	currentSong, err := c.CurrentSong()
	c.Close()
	config.Log(err)

	q.setPageInfo(page, currentSong["Pos"])
	return q.getQueueBounds()
}

func (q Queue) getQueueBounds() (downLimit, upLimt int) {
	downLimit = int((q.CurrentPage - 1) * ClientQueueLimit)
	upLimt = downLimit + ClientQueueLimit

	if upLimt > int(q.Length) {
		upLimt = int(q.Length)
	}
	return
}

// sets the value of cuerrnt page
func (q *Queue) setPageInfo(page, currentSong string) {
	csPos, _ := strconv.Atoi(currentSong)
	q.CurrentSongPage = uint(csPos/ClientQueueLimit + 1)
	q.LastPage = uint(q.Length/ClientQueueLimit + 1)

	parsedPage, err := strconv.Atoi(page)
	if err != nil || parsedPage < 1 {
		q.CurrentPage = q.CurrentSongPage
		return
	}

	if uint(parsedPage) > q.LastPage {
		q.CurrentPage = q.LastPage
		return
	}

	q.CurrentPage = uint(parsedPage)
}
