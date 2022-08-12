package playback

import (
	"fmt"
	"strconv"
	"time"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/song"
)

var err error

type Commands map[string]func(c *config.Connection) error

var commands = Commands{
	"clear":        clearQueue,
	"consume":      consume,
	"next":         nextSong,
	"previous":     prevSong,
	"random":       random,
	"repeat":       repeat,
	"single":       single,
	"stop":         stop,
	"toggle":       toggle,
	"play":         play,
	"seekForward":  seekForward,
	"seekBackward": seekBackward,
	"volumeDown":   volumeDown,
	"volumeUp":     volumeUp,
}

func Command(c config.Connection, cmd string, data int) error {
	c.Connect()
	defer c.Close()
	if command, ok := commands[cmd]; ok {
		return command(&c)
	}

	switch cmd {
	case "seek":
		err = seek(&c, data)
	case "changeVolume":
		err = changeVolume(&c, data)
	default:
		return fmt.Errorf("invalid command")
	}

	return err
}

// toggles the state between play and pause
func toggle(c *config.Connection) error {
	status := song.GetStatus(*c)
	if status["state"] == "play" {
		err = c.Client.Pause(true)
		return err
	}
	err = c.Client.Play(-1)
	return err
}

// play starts playing the current song in queue
func play(c *config.Connection) error {
	err = c.Client.Play(-1)
	config.Log(err)
	return err
}

// stops the current queue
func stop(c *config.Connection) error {
	err = c.Client.Stop()
	config.Log(err)
	return err
}

// goes to next song in queue
func nextSong(c *config.Connection) error {
	err = c.Client.Next()
	config.Log(err)
	return err
}

// goes to Previous song in queue
func prevSong(c *config.Connection) error {
	err = c.Client.Previous()
	config.Log(err)
	return err
}

// seeks t second in current song
func seek(c *config.Connection, t int) error {
	seekDuration, err := time.ParseDuration(fmt.Sprintf("%ds", t))
	config.Log(err)
	err = c.Client.SeekCur(seekDuration, false)
	config.Log(err)
	return err
}

// seeks 15 seconds forward in current song
func seekForward(c *config.Connection) error {
	err = c.Client.SeekCur(15*time.Second, true)
	config.Log(err)
	return err
}

// seeks 15 seconds backward in current song
func seekBackward(c *config.Connection) error {
	err = c.Client.SeekCur(-15*time.Second, true)
	config.Log(err)
	return err
}

// toggles the Single state
func single(c *config.Connection) error {
	status := song.GetStatus(*c)
	single := false
	if status["single"] == "0" {
		single = true
	}
	err = c.Client.Single(single)
	config.Log(err)
	return err
}

// toggles the Repeat state
func repeat(c *config.Connection) error {
	status := song.GetStatus(*c)
	repeat := false
	if status["repeat"] == "0" {
		repeat = true
	}
	err = c.Client.Repeat(repeat)
	config.Log(err)
	return err
}

// toggles the Consume state
func consume(c *config.Connection) error {
	status := song.GetStatus(*c)
	consume := false
	if status["consume"] == "0" {
		consume = true
	}
	err = c.Client.Consume(consume)
	config.Log(err)
	return err
}

// toggles the random state
func random(c *config.Connection) error {
	status := song.GetStatus(*c)
	repeat := false
	if status["random"] == "0" {
		repeat = true
	}
	err = c.Client.Random(repeat)
	config.Log(err)
	return err
}

// sets the volume
func changeVolume(c *config.Connection, volume int) error {
	if volume > 100 {
		volume = 100
	}
	if volume < 0 {
		volume = 0
	}

	err = c.Client.SetVolume(volume)
	config.Log(err)
	return err
}

// returns the volume
func getVolume(c *config.Connection) int {
	volAttr, err := c.Client.Command("getvol").Attrs()
	config.Log(err)
	volume, _ := strconv.Atoi(volAttr["volume"])
	return volume
}

// Increases the  Volume by 5%
func volumeUp(c *config.Connection) error {
	volume := getVolume(c)
	err = changeVolume(c, volume+5)
	config.Log(err)
	return err
}

// Decreases the Volume by 5%
func volumeDown(c *config.Connection) error {
	volume := getVolume(c)
	err = changeVolume(c, volume-5)
	config.Log(err)
	return err
}

// clears the current queue
func clearQueue(c *config.Connection) error {
	err = c.Client.Clear()
	config.Log(err)
	return err
}
