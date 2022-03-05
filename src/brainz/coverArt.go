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

func Fetch(url string) (body []byte, err error) {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timeout, http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "ava-mpd/1.2.0 ( smohsenk@gmail.com )")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		config.Log(err)
		return body, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNotFound:
		return body, fmt.Errorf("%d : couldn't find resource for %s", http.StatusNotFound, url)
	case http.StatusBadRequest:
		return body, fmt.Errorf("%d : couldn't find resource for %s", http.StatusBadRequest, url)
	case http.StatusServiceUnavailable:
		return body, fmt.Errorf("%d : rate limited %s", http.StatusServiceUnavailable, url)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func GetCoverURL(releaseID string) ([]byte, error) {
	res, err := Fetch(CAURL + releaseID)
	if err != nil {
		return nil, err
	}

	url := parseImageURL(res)
	if url == "" {
		return nil, fmt.Errorf("couldn't find cover")
	}
	cover, err := Fetch(url)

	return cover, err
}

func parseImageURL(b []byte) (url string) {
	resp := CoverArtResponse{}
	if err := json.Unmarshal(b, &resp); err != nil {
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
