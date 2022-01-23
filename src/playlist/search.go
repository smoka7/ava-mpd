package playlist

import (
	"errors"

	"github.com/smoka7/ava/src/config"
)

type SearchResult map[string][]File

//search for songs in the server
func SearchServer(c *config.Connection, term ...string) (result SearchResult) {
	filter := term[0]
	err := validFilter(filter)
	config.Log(err)
	result = make(map[string][]File)
	query, err := c.Client.Search(term...)
	config.Log(err)
	if len(query) >= 400 {
		query = query[:400]
	}
	for i := 0; i < len(query); i++ {
		if _, ok := result[query[i][filter]]; !ok {
			result[query[i][filter]] = make([]File, 0)
		}
		result[query[i][filter]] = append(result[query[i][filter]], newFile(query[i]))
	}
	return
}

func validFilter(filter string) error {
	v := map[string]bool{
		"file": true, "Artist": true, "Album": true,
		"Genre": true, "Date": true, "Title": true,
	}
	if _, ok := v[filter]; !ok {
		return errors.New("valid filters are : File,Artist,Album,Genre,Date and Title")
	}
	return nil
}
