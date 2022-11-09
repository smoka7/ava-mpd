package brainz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/smoka7/ava/src/config"
)

// api endpoint of coverArts
const CAURL = "https://coverartarchive.org/release/"

type CoverArtResponse struct {
	Images  []images `json:"images"`
	Release string   `json:"release"`
}

type images struct {
	Image string `json:"image"`
	Front bool   `json:"front"`
}

func fetch(url string) ([]byte, error) {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timeout, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ava-mpd/1.2.0 ( smohsenk@gmail.com )")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		config.Log(err)
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("%d : couldn't find resource for %s", http.StatusNotFound, url)
	case http.StatusBadRequest:
		return nil, fmt.Errorf("%d : couldn't find resource for %s", http.StatusBadRequest, url)
	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("%d : rate limited %s", http.StatusServiceUnavailable, url)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// gets the cover for release releaseID from musicbrainz
func GetCover(releaseID string) ([]byte, error) {
	res, err := fetch(CAURL + releaseID)
	if err != nil {
		return nil, err
	}

	url := getCoverURL(res)
	if url == "" {
		return nil, fmt.Errorf("couldn't find cover")
	}
	cover, err := fetch(url)

	return cover, err
}

// parses the cover art response and returns the cover url
func getCoverURL(body []byte) string {
	var resp CoverArtResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		config.Log(err)
		return ""
	}

	for _, image := range resp.Images {
		if image.Front {
			return image.Image
		}
	}

	return ""
}
