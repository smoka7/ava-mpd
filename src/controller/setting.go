package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

type (
	settingCommand  string
	settingCommands map[settingCommand]func(int)
)

type SettingRequest struct {
	Command settingCommand `json:"command"`
	Data    settingData    `json:"data"`
}

type settingData struct {
	Value int
}

type SettingsResponse struct {
	Outputs          config.Outputs
	DatabaseStats    mpd.Attrs
	ReplayGain       string
	DownloadCoverArt bool
}

var (
	crossfade      settingCommand = "crossfade"
	download       settingCommand = "download"
	mixrampdb      settingCommand = "mixrampdb"
	enableOutput   settingCommand = "enableOutput"
	disableOutput  settingCommand = "disableOutput"
	deleteCache    settingCommand = "deleteCache"
	setGain        settingCommand = "setGain"
	updateDatabase settingCommand = "updateDatabase"
)

func (c Mpd) Settings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request SettingRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		c.setingcommand(request)

	case http.MethodGet:
		response := SettingsResponse{
			DatabaseStats:    c.Client.DatabaseStats(),
			ReplayGain:       c.Client.GetReplayGain(),
			Outputs:          c.Client.ListOutputs(),
			DownloadCoverArt: c.Client.DownloadCoverArt,
		}
		err := json.NewEncoder(w).Encode(response)
		config.Log(err)
	}
}

func (c Mpd) setingcommand(request SettingRequest) {
	c.Client.Connect()
	defer c.Client.Close()

	cmd := settingCommands{
		crossfade: func(value int) {
			c.Client.ChangeCrossfade(value)
		},
		download: func(value int) {
			c.Client.ToggleDownloadCover()
		},
		mixrampdb: func(value int) {
			c.Client.ChangeMixRampdb(value)
		},
		enableOutput: func(value int) {
			c.Client.ToggleOutput(value, true)
		},
		disableOutput: func(value int) {
			c.Client.ToggleOutput(value, false)
		},
		deleteCache: func(value int) {
			config.DeleteCache()
		},
		setGain: func(value int) {
			c.Client.ChangeReplayGain(value)
		},
		updateDatabase: func(value int) {
			c.Client.UpdateDatabase()
		},
	}

	if action, ok := cmd[request.Command]; ok {
		action(request.Data.Value)
	}
}
