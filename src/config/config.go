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
	AppPort  string      //port of current app
	Client   *mpd.Client //connected client
}

//maximm log size in bytes
const MaxLogSize = 20000

const LogFilePath = "/ava-mpd/ava.log"

var err error

//Reads the MPD server connection from environment values
func (c *Connection) ReadEnv() {
	host := os.Getenv("MPD_HOST")
	port := os.Getenv("MPD_PORT")
	if host != "" && port != "" {
		c.Address = host + ":" + port
		c.Password = os.Getenv("MPD_PASSWORD")
	}
}

//parse the MPD server connection from bin flag
func (c *Connection) ReadFromFlags() {
	var address string
	flag.StringVar(&address, "address", "", "address of mpd server host:port")
	flag.StringVar(&c.Password, "password", c.Password, "password of mpd server")
	flag.StringVar(&c.AppPort, "port", "3001", "The port to run this app on it")
	flag.Parse()
	if address != "" {
		c.Address = address
	}
}

//connects to server
func (c *Connection) Connect() {
	c.Client, err = mpd.DialAuthenticated("tcp", c.Address, c.Password)
	Log(err)
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
		return
	}
	logFile, e := os.OpenFile(cache+LogFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if e != nil {
		log.Println(e)
		return
	}
	stat, _ := logFile.Stat()
	if stat.Size() > MaxLogSize {
		logFile.Close()
		logFile, e = os.OpenFile(cache+LogFilePath, os.O_TRUNC|os.O_WRONLY, 0600)
		if e != nil {
			log.Println(e)
			return
		}
	}
	Log := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	Log.Println(err)
	logFile.Close()
}

//updates the MPD server database
func UpdateDatabase(c Connection) {
	_, err := c.Client.Update("")
	Log(err)
}

//returns the MPD database stats
func DatabaseStats(c Connection) (stats mpd.Attrs) {
	stats, err := c.Client.Stats()
	Log(err)
	return
}

//returns the mpd outputs
func ListOutputs(c Connection) (stats []mpd.Attrs) {
	stats, err := c.Client.ListOutputs()
	Log(err)
	return
}

//enables the output
func EnableOutput(c Connection, id int) {
	err := c.Client.EnableOutput(id)
	Log(err)
}

//Disables the output
func DisableOutput(c Connection, id int) {
	err := c.Client.DisableOutput(id)
	Log(err)
}

//sets the crossFade
func ChangeCrossfade(c Connection, second int) {
	err := c.Client.Command("crossfade %d", second).OK()
	Log(err)
}

//sets the Gain status
func ChangeReplayGain(c Connection, id int) {
	modes := map[int]string{
		0: "off",
		1: "track",
		2: "album",
		3: "auto",
	}
	err := c.Client.Command("replay_gain_mode %s", modes[id]).OK()
	Log(err)
}

//returns the Gain status
func GetReplayGain(c Connection) (status mpd.Attrs) {
	status, err = c.Client.Command("replay_gain_status").Attrs()
	Log(err)
	return
}
