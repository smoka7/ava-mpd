package config

import (
	"flag"
	"log"
	"os"

	"github.com/fhs/gompd/v2/mpd"
)

type Connection struct {
	Address  string      //address of mpd server
	Password string      //password of mpd server
	Client   *mpd.Client //connected client
}

var err error

//Reads the MPD server connection from environment values
func (c *Connection) ReadEnv() {
	c.Address = os.Getenv("MPD_HOST") + ":" + os.Getenv("MPD_PORT")
	c.Password = os.Getenv("MPD_PASSWORD")
}

//parse the MPD server connection from bin flag
func (c *Connection) ReadFromFlags() {
	var host, port string
	flag.StringVar(&host, "host", host, "host of mpd server")
	flag.StringVar(&port, "port", port, "port of mpd server")
	flag.StringVar(&c.Password, "password", c.Password, "password of mpd server")
	flag.Parse()
	if host != "" && port != "" {
		c.Address = host + ":" + port
	}
}

//connects to server
func (c *Connection) Connect() {
	if c.Address == "" {
		log.Println("enter host and port of server")
		os.Exit(1)
	}
	c.Client, err = mpd.DialAuthenticated("tcp", c.Address, c.Password)
	if err != nil {
		Log(err)
	}
}

//closes the connection to server
func (c *Connection) Close() {
	c.Client.Close()
}

//deletes the albumArts cache
func DeleteCache() {
	cache, err := os.UserCacheDir()
	Log(err)
	coverArtDir := cache + "/ava-mpd/coverart"
	err = os.RemoveAll(coverArtDir)
	Log(err)
}

// checks log and saves it to cache folder
func Log(err error) {
	if err == nil {
		return
	}
	cache, e := os.UserCacheDir()
	if e != nil {
		log.Println(e)
	}
	logFile, e := os.OpenFile(cache+"/ava-mpd/ava.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if e != nil {
		log.Println(e)
	}
	defer logFile.Close()
	Log := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	Log.Println(err)
}

//updates the MPD server database
func UpdateDatabase(c Connection) {
	c.Connect()
	_, err := c.Client.Update("")
	Log(err)
	c.Close()
}

//returns the MPD database stats
func DatabaseStats(c Connection) (stats mpd.Attrs) {
	c.Connect()
	stats, err := c.Client.Stats()
	Log(err)
	c.Close()
	return
}

//returns the mpd outputs
func ListOutputs(c Connection) (stats []mpd.Attrs) {
	c.Connect()
	stats, err := c.Client.ListOutputs()
	Log(err)
	c.Close()
	return
}

//enables the output
func EnableOutput(c Connection, id int) {
	c.Connect()
	err := c.Client.EnableOutput(id)
	Log(err)
	c.Close()
}

//Disables the output
func DisableOutput(c Connection, id int) {
	c.Connect()
	err := c.Client.DisableOutput(id)
	Log(err)
	c.Close()
}

//sets the crossFade
func ChangeCrossfade(c Connection, second int) {
	c.Connect()
	err := c.Client.Command("crossfade %d", second).OK()
	Log(err)
	c.Close()
}

//sets the Gain status
func ChangeReplayGain(c Connection, id int) {
	c.Connect()
	modes := map[int]string{
		0: "off",
		1: "track",
		2: "album",
		3: "auto",
	}
	err := c.Client.Command("replay_gain_mode %s", modes[id]).OK()
	Log(err)
	c.Close()
}

//returns the Gain status
func GetReplayGain(c Connection) (status mpd.Attrs) {
	c.Connect()
	status, err = c.Client.Command("replay_gain_status").Attrs()
	Log(err)
	c.Close()
	return
}
