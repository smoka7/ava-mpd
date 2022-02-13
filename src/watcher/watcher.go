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

type eventMessage struct {
	Subsystem string
}

func NewClient(c config.Connection) Mpd {
	return Mpd{Client: c}
}

// sends the mpd events to client's Socket
func (m *Mpd) Serve(ws *websocket.Conn) {
	event := make(chan string, 4)
	message := eventMessage{}
	eventWatcher(m.Client, event)
	for {
		message.Subsystem = <-event
		err := websocket.JSON.Send(ws, <-event)
		config.Log(err)
	}
}

// connect to mpd server and watches for the subsystem change
func eventWatcher(c config.Connection, event chan string) {
	watcher, err := mpd.NewWatcher("tcp", c.Address, c.Password)
	config.Log(err)
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
func (m *Mpd) RecordPlayCount() {
	client := m.Client
	watcher, err := mpd.NewWatcher("tcp", client.Address, client.Password)
	config.Log(err)
	defer watcher.Close()
	go func() {
		for err := range watcher.Error {
			config.Log(err)
		}
	}()
	func() {
		client.Connect()
		status := song.GetStatus(&client)["state"]
		lastSong := song.GetCurrentSong(client)
		if status == "play" {
			song.IncrementPCount(&client)
		}
		client.Close()
		for subsystem := range watcher.Event {
			if subsystem == "player" {
				client.Connect()
				currSong := song.GetCurrentSong(client)
				currStatus := song.GetStatus(&client)["state"]
				if lastSong["file"] != currSong["file"] &&
					currStatus == "play" {
					song.IncrementPCount(&client)
				}
				client.Close()
				lastSong = currSong
			}
		}
	}()
}
