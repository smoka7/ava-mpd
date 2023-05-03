package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/smoka7/ava/src/config"
	"github.com/smoka7/ava/src/playlist"
)

type SearchRequest struct {
	Command string   `json:"command"`
	Terms   []string `json:"terms"`
}

type SearchResponse struct {
	Songs playlist.SearchResult
	Error error
}

func (c *Mpd) SearchServer(w http.ResponseWriter, r *http.Request) {
	c.Connect()
	defer c.Close()
	if r.Method == http.MethodPost {
		var request SearchRequest
		var files playlist.SearchResult
		err := json.NewDecoder(r.Body).Decode(&request)
		config.Log(err)
		switch request.Command {
		case "search":
			files, err = playlist.SearchServer(c.Connection, request.Terms...)
		case "searchadd":
			err = playlist.SearchAdd(c.Connection, request.Terms...)
		case "searchplay":
			err = playlist.SearchPlay(c.Connection, request.Terms...)
		default:
			err = errors.New("chose a command")
		}
		err = json.NewEncoder(w).Encode(SearchResponse{
			Songs: files,
			Error: err,
		})
		config.Log(err)
		return
	}
}
