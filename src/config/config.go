package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/fhs/gompd/v2/mpd"
)

type Connection struct {
	DownloadCoverArt bool        `json:"download_cover_art"` // whether to Download missing coverart from musicbrainz
	Address          string      `json:"address"`            // address of mpd server
	Password         string      `json:"password,omitempty"` // password of mpd server
	AppPort          string      `json:"app_port"`           // port of current app
	Client           *mpd.Client `json:"-"`                  // connected client
}

const (
	MaxLogSize     = 20000 // maximm log size in bytes
	LogFilePath    = "/ava-mpd/ava.log"
	ConfigFilePath = "/ava-mpd/"
	ConfigFileName = "ava.json"
)

var err error

// Reads the MPD server connection from environment values
func (c *Connection) readEnv() {
	host := os.Getenv("MPD_HOST")
	port := os.Getenv("MPD_PORT")
	if host != "" && port != "" {
		c.Address = host + ":" + port
	}
}

// parse the MPD server connection from bin flag
func (c *Connection) ReadConfigs() {
	var configPath string
	var flagConfig Connection

	flag.StringVar(&flagConfig.Address, "address", "", "address of mpd server host:port")
	flag.StringVar(&flagConfig.Password, "password", c.Password, "password of mpd server")
	flag.StringVar(&flagConfig.AppPort, "port", "3001", "The port to run this app on it defaults to 3001")
	flag.StringVar(&configPath, "config", "", "absolute path of config file")
	flag.Parse()

	c.readEnv()

	if configPath != "" {
		err := c.loadConfig(configPath)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	if c.Address == "" {
		c.Address = flagConfig.Address
	}

	if c.Password == "" {
		c.Password = flagConfig.Password
	}

	if c.AppPort == "" {
		c.AppPort = flagConfig.AppPort
	}

	if c.Address == "" {
		log.Println("mpd server address is empty")
		os.Exit(1)
	}
}

// saves the app configurations to file
func (c *Connection) SaveConfig() {
	bytes, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		Log(err)
		return
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		Log(err)
		return
	}

	configPath := configDir + ConfigFilePath
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		err = os.Mkdir(configPath, 0600)
		if err != nil {
			Log(err)
			return
		}
	}

	err = os.WriteFile(configPath+ConfigFileName, bytes, 0600)
	Log(err)
}

func (c *Connection) loadConfig(configPath string) error {
	// when config path is omited use the default one
	if configPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			Log(err)
			return err
		}
		configPath = configDir + ConfigFilePath + ConfigFileName
	}

	// if config file doesn't exist let it go
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}

	return nil
}

// connects to server
func (c *Connection) Connect() {
	c.Client, err = mpd.DialAuthenticated("tcp", c.Address, c.Password)
	Log(err)
}

// closes the connection to server
func (c *Connection) Close() {
	c.Client.Close()
}

// deletes the albumArts cache
func DeleteCache() {
	cache, err := os.UserCacheDir()
	Log(err)
	coverArtDir := cache + "/ava-mpd/coverart"
	err = os.RemoveAll(coverArtDir)
	Log(err)
}

//  checks log and saves it to cache folder
func Log(err error) {
	if err == nil {
		return
	}
	cache, e := os.UserCacheDir()
	if e != nil {
		log.Println(e)
		return
	}
	logFile, e := os.OpenFile(cache+LogFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if e != nil {
		log.Println(e)
		return
	}
	stat, _ := logFile.Stat()
	if stat.Size() > MaxLogSize {
		logFile.Close()
		logFile, e = os.OpenFile(cache+LogFilePath, os.O_TRUNC|os.O_WRONLY, 0o600)
		if e != nil {
			log.Println(e)
			return
		}
	}
	Log := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	Log.Println(err)
	logFile.Close()
}

// updates the MPD server database
func (c *Connection) UpdateDatabase() {
	_, err := c.Client.Update("")
	Log(err)
}

// returns the MPD database stats
func (c Connection) DatabaseStats() (stats mpd.Attrs) {
	c.Connect()
	defer c.Close()
	stats, err := c.Client.Stats()
	Log(err)

	return
}

// toggles wheter to download missing covers
func (c *Connection) ToggleDownloadCover() {
	c.DownloadCoverArt = !c.DownloadCoverArt
	c.SaveConfig()
}

// returns the mpd outputs
func (c Connection) ListOutputs() (stats []mpd.Attrs) {
	c.Connect()
	defer c.Close()
	stats, err := c.Client.ListOutputs()
	Log(err)
	return
}

// enables the output
func (c *Connection) EnableOutput(id int) {
	err := c.Client.EnableOutput(id)
	Log(err)
}

// Disables the output
func (c *Connection) DisableOutput(id int) {
	err := c.Client.DisableOutput(id)
	Log(err)
}

// sets the crossFade
func (c *Connection) ChangeCrossfade(second int) {
	err := c.Client.Command("crossfade %d", second).OK()
	Log(err)
}

// sets the crossFade
func (c *Connection) ChangeMixRampdb(second int) {
	err := c.Client.Command("mixrampdb %d", second).OK()
	Log(err)
}

// sets the Gain status
func (c *Connection) ChangeReplayGain(id int) {
	modes := map[int]string{
		0: "off",
		1: "track",
		2: "album",
		3: "auto",
	}
	if mode, ok := modes[id]; ok {
		err := c.Client.Command("replay_gain_mode %s", mode).OK()
		Log(err)
	}
}

// returns the Gain status
func (c Connection) GetReplayGain() string {
	c.Connect()
	defer c.Close()
	status, err := c.Client.Command("replay_gain_status").Attrs()
	if err != nil {
		return ""
	}
	Log(err)
	return status["replay_gain_mode"]
}
