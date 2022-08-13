package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/smoka7/ava/src/config"
)

type SettingRequest struct {
	Command string      `json:"command"`
	Data    SettingData `json:"data"`
}

type SettingData struct {
	Value int
}

type SettingsResponse struct {
	Outputs          []mpd.Attrs
	DatabaseStats    mpd.Attrs
	ReplayGain       string
	DownloadCoverArt bool
}

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
	switch request.Command {
	case "crossfade":
		c.Client.ChangeCrossfade(request.Data.Value)
	case "download":
		c.Client.ToggleDownloadCover()
	case "mixrampdb":
		c.Client.ChangeMixRampdb(request.Data.Value)
	case "enableOutput":
		c.Client.EnableOutput(request.Data.Value)
	case "disableOutput":
		c.Client.DisableOutput(request.Data.Value)
	case "deleteCache":
		config.DeleteCache()
	case "setGain":
		c.Client.ChangeReplayGain(request.Data.Value)
	case "updateDatabase":
		c.Client.UpdateDatabase()
	}
}
