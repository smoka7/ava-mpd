package playback

import (
	"fmt"
	"time"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/song"
)

var err error

//toggles the state between play and pause
func Toggle(c *config.Connection) {
	status := song.GetStatus(c)
	if status["state"] == "play" {
		c.Client.Pause(true)
		return
	}
	c.Client.Play(-1)
}

//stops the current queue
func Stop(c *config.Connection) {
	err = c.Client.Stop()
	config.Log(err)
}

//goes to next song in queue
func NextSong(c *config.Connection) {
	err = c.Client.Next()
	config.Log(err)
}

//goes to Previous song in queue
func PrevSong(c *config.Connection) {
	err = c.Client.Previous()
	config.Log(err)
}

//seeks t second in current song
func Seek(c *config.Connection, t int) {
	seekDuration, err := time.ParseDuration(fmt.Sprintf("%ds", t))
	config.Log(err)
	err = c.Client.SeekCur(seekDuration, false)
	config.Log(err)
}

//toggles the Single state
func Single(c *config.Connection) {
	status := song.GetStatus(c)
	single := false
	if status["single"] == "0" {
		single = true
	}
	err = c.Client.Single(single)
	config.Log(err)
}

//toggles the Repeat state
func Repeat(c *config.Connection) {
	status := song.GetStatus(c)
	repeat := false
	if status["repeat"] == "0" {
		repeat = true
	}
	err = c.Client.Repeat(repeat)
	config.Log(err)
}

//toggles the Consume state
func Consume(c *config.Connection) {
	status := song.GetStatus(c)
	consume := false
	if status["consume"] == "0" {
		consume = true
	}
	err = c.Client.Consume(consume)
	config.Log(err)
}

//toggles the random state
func Random(c *config.Connection) {
	status := song.GetStatus(c)
	repeat := false
	if status["random"] == "0" {
		repeat = true
	}
	err = c.Client.Random(repeat)
	config.Log(err)
}

//sets the volume
func ChangeVolume(c *config.Connection, volume int) {
	err = c.Client.SetVolume(volume)
	config.Log(err)
}

//clears the current queue
func ClearQueue(c *config.Connection) {
	err = c.Client.Clear()
	config.Log(err)
}
