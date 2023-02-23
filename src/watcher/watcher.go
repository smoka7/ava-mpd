package watcher

import (
	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/song"
	"golang.org/x/net/websocket"
)

type Mpd struct {
	Client config.Connection
}

type status struct {
	state, path string
}

type eventMessage struct {
	Subsystem string
}

func NewClient(c config.Connection) Mpd {
	return Mpd{Client: c}
}

// sends the mpd events to client's Socket
func (m Mpd) Serve(ws *websocket.Conn) {
	event := make(chan string, 4)
	message := eventMessage{}
	eventWatcher(m.Client, event)
	for {
		message.Subsystem = <-event
		err := websocket.JSON.Send(ws, message)
		config.Log(err)
	}
}

// connect to mpd server and watches for the subsystem change
func eventWatcher(c config.Connection, event chan string) {
	watcher, err := mpd.NewWatcher("tcp", c.Address, c.Password)
	config.LogAndExit(err)

	go func() {
		for err := range watcher.Error {
			config.Log(err)
		}
	}()
	go func() {
		for subsystem := range watcher.Event {
			event <- subsystem
		}
	}()
}

// saves the each songs play count
func (m Mpd) RecordPlayCount() {
	client := m.Client
	watcher, err := mpd.NewWatcher("tcp", client.Address, client.Password)
	config.LogAndExit(err)

	defer watcher.Close()

	go func() {
		for err := range watcher.Error {
			config.Log(err)
		}
	}()

	status := newStatus()
	status.compareStatus(client)
	for subsystem := range watcher.Event {
		if subsystem != "player" {
			continue
		}
		status.compareStatus(client)
	}
}

func newStatus() status {
	return status{
		path:  "",
		state: "",
	}
}

// compareStatus compares current and last playing status of server
// and increments playing count accordingly
func (s *status) compareStatus(c config.Connection) {
	c.Connect()
	defer c.Close()

	s.state = song.GetStatus(c)["state"]
	currentSong, err := c.Client.CurrentSong()
	config.Log(err)

	if s.path != currentSong["file"] &&
		s.state == "play" {
		song.IncrementPCount(c, currentSong["file"])
	}
	s.path = currentSong["file"]
}
