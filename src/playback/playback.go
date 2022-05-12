package playback

import (
	"fmt"
	"time"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/song"
)

var err error

func Command(c config.Connection, cmd string, data int) {
	c.Connect()
	defer c.Close()
	switch cmd {
	case "next":
		nextSong(&c)
	case "previous":
		prevSong(&c)
	case "toggle":
		toggle(&c)
	case "stop":
		stop(&c)
	case "single":
		single(&c)
	case "repeat":
		repeat(&c)
	case "random":
		random(&c)
	case "consume":
		consume(&c)
	case "clear":
		clearQueue(&c)
	case "seek":
		seek(&c, data)
	case "changeVolume":
		changeVolume(&c, data)
	}
}

//  toggles the state between play and pause
func toggle(c *config.Connection) {
	status := song.GetStatus(*c)
	if status["state"] == "play" {
		c.Client.Pause(true)
		return
	}
	c.Client.Play(-1)
}

// stops the current queue
func stop(c *config.Connection) {
	err = c.Client.Stop()
	config.Log(err)
}

// goes to next song in queue
func nextSong(c *config.Connection) {
	err = c.Client.Next()
	config.Log(err)
}

// goes to Previous song in queue
func prevSong(c *config.Connection) {
	err = c.Client.Previous()
	config.Log(err)
}

// seeks t second in current song
func seek(c *config.Connection, t int) {
	seekDuration, err := time.ParseDuration(fmt.Sprintf("%ds", t))
	config.Log(err)
	err = c.Client.SeekCur(seekDuration, false)
	config.Log(err)
}

// toggles the Single state
func single(c *config.Connection) {
	status := song.GetStatus(*c)
	single := false
	if status["single"] == "0" {
		single = true
	}
	err = c.Client.Single(single)
	config.Log(err)
}

// toggles the Repeat state
func repeat(c *config.Connection) {
	status := song.GetStatus(*c)
	repeat := false
	if status["repeat"] == "0" {
		repeat = true
	}
	err = c.Client.Repeat(repeat)
	config.Log(err)
}

// toggles the Consume state
func consume(c *config.Connection) {
	status := song.GetStatus(*c)
	consume := false
	if status["consume"] == "0" {
		consume = true
	}
	err = c.Client.Consume(consume)
	config.Log(err)
}

// toggles the random state
func random(c *config.Connection) {
	status := song.GetStatus(*c)
	repeat := false
	if status["random"] == "0" {
		repeat = true
	}
	err = c.Client.Random(repeat)
	config.Log(err)
}

// sets the volume
func changeVolume(c *config.Connection, volume int) {
	err = c.Client.SetVolume(volume)
	config.Log(err)
}

// clears the current queue
func clearQueue(c *config.Connection) {
	err = c.Client.Clear()
	config.Log(err)
}
