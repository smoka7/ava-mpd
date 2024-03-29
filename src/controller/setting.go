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
	ReplayGain       string         `json:"ReplayGain"`
	DatabaseStats    mpd.Attrs      `json:"DatabaseStats"`
	Outputs          config.Outputs `json:"Outputs"`
	DownloadCoverArt bool           `json:"DownloadCoverArt"`
}

const (
	crossfade      settingCommand = "crossfade"
	download       settingCommand = "download"
	mixrampdb      settingCommand = "mixrampdb"
	toggleOutput   settingCommand = "toggleOutput"
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
			DatabaseStats:    c.DatabaseStats(),
			ReplayGain:       c.GetReplayGain(),
			Outputs:          c.GetOutputs(),
			DownloadCoverArt: c.DownloadCoverArt,
		}
		err := json.NewEncoder(w).Encode(response)
		config.Log(err)
	}
}

func (c Mpd) setingcommand(request SettingRequest) {
	c.Connect()
	defer c.Close()

	cmd := settingCommands{
		download:       func(_ int) { c.ToggleDownloadCover() },
		updateDatabase: func(_ int) { c.UpdateDatabase() },
		deleteCache:    func(_ int) { config.DeleteCache() },
		toggleOutput:   c.ToggleOutput,
		crossfade:      c.ChangeCrossfade,
		mixrampdb:      c.ChangeMixRampdb,
		setGain:        c.ChangeReplayGain,
	}

	if action, ok := cmd[request.Command]; ok {
		action(request.Data.Value)
	}
}
