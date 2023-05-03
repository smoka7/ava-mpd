package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/fhs/gompd/v2/mpd"
)

type Connection struct {
	*mpd.Client      `json:"-"` // connected client
	Address          string     `json:"address"`            // address of mpd server
	Password         string     `json:"password,omitempty"` // password of mpd server
	AppPort          string     `json:"app_port"`           // port of current app
	DownloadCoverArt bool       `json:"download_cover_art"` // whether to Download missing coverart from musicbrainz
}

const (
	MaxLogSize     = 20000 // maximm log size in bytes
	LogFilePath    = "/ava-mpd/ava.log"
	ConfigFilePath = "/ava-mpd/"
	ConfigFileName = "ava.json"
)

// return a connection with default values.
func newDefaultConfig() Connection {
	return Connection{
		Address:          "localhost:6600",
		AppPort:          "3001",
		DownloadCoverArt: false,
		Password:         "",
	}
}

// copies the config from to connection to struct and ignores the empty values.
func (c *Connection) copyConfig(newConfig Connection) {
	if newConfig.Address != "" {
		c.Address = newConfig.Address
	}

	if newConfig.Password != "" {
		c.Password = newConfig.Password
	}

	if newConfig.AppPort != "" {
		c.AppPort = newConfig.AppPort
	}
	c.DownloadCoverArt = newConfig.DownloadCoverArt
}

// Reads the MPD server connection from environment values.
func (c *Connection) readEnv() {
	host := os.Getenv("MPD_HOST")
	port := os.Getenv("MPD_PORT")
	if host != "" && port != "" {
		c.Address = host + ":" + port
	}
}

func (c *Connection) parseFlags() (configPath string) {
	var flags Connection

	flag.StringVar(&flags.Address, "address", "", "address of mpd server host:port")
	flag.StringVar(&flags.Password, "password", "", "password of mpd server")
	flag.StringVar(&flags.AppPort, "port", "", "The port to run this app on it defaults to 3001")
	flag.StringVar(&configPath, "config", "", "absolute path of config file")

	flag.Parse()
	c.copyConfig(flags)
	return
}

// parse the MPD server connection from bin flag.
func ReadConfigs() Connection {
	c := newDefaultConfig()

	c.readEnv()

	configPath := c.parseFlags()

	c.loadConfigFile(configPath)
	return c
}

// saves the app configurations to file.
func (c Connection) SaveConfig() {
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

func (c *Connection) loadConfigFile(configPath string) {
	// when config path is omitted use the default one
	var config Connection
	if configPath == "" {
		configDir, err := os.UserConfigDir()
		LogAndExit(err)

		configPath = configDir + ConfigFilePath + ConfigFileName
	}

	// if config file doesn't exist let it go
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	bytes, err := os.ReadFile(configPath)
	LogAndExit(err)

	err = json.Unmarshal(bytes, &config)
	LogAndExit(err)

	c.copyConfig(config)
}

// connects to server.
func (c *Connection) Connect() {
	client, err := mpd.DialAuthenticated("tcp", c.Address, c.Password)
	Log(err)
	c.Client = client
}

// closes the connection to server.
func (c *Connection) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}

// deletes the albumArts cache.
func DeleteCache() {
	cache, err := os.UserCacheDir()
	Log(err)
	coverArtDir := cache + "/ava-mpd/coverart"
	err = os.RemoveAll(coverArtDir)
	Log(err)
}

// checks log and saves it to cache folder.
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

// for none nil errors logs to file and then exits.
func LogAndExit(err error) {
	if err != nil {
		Log(err)
		log.Fatalln(err)
	}
}

// updates the MPD server database.
func (c Connection) UpdateDatabase() {
	_, err := c.Update("")
	Log(err)
}

// returns the MPD database stats.
func (c Connection) DatabaseStats() (stats mpd.Attrs) {
	c.Connect()
	defer c.Close()
	stats, err := c.Stats()
	Log(err)

	return
}

// toggles whether to download missing covers.
func (c *Connection) ToggleDownloadCover() {
	c.DownloadCoverArt = !c.DownloadCoverArt
	c.SaveConfig()
}

// sets the crossFade.
func (c Connection) ChangeCrossfade(second int) {
	err := c.Command("crossfade %d", second).OK()
	Log(err)
}

// sets the crossFade.
func (c Connection) ChangeMixRampdb(second int) {
	err := c.Command("mixrampdb %d", second).OK()
	Log(err)
}

// sets the Gain status.
func (c Connection) ChangeReplayGain(id int) {
	modes := map[int]string{
		0: "off",
		1: "track",
		2: "album",
		3: "auto",
	}
	if mode, ok := modes[id]; ok {
		err := c.Command("replay_gain_mode %s", mode).OK()
		Log(err)
	}
}

// returns the Gain status.
func (c Connection) GetReplayGain() string {
	c.Connect()
	defer c.Close()
	status, err := c.Command("replay_gain_status").Attrs()
	if err != nil {
		return ""
	}
	Log(err)
	return status["replay_gain_mode"]
}
