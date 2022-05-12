package playlist

import (
	"errors"

	"github.com/smoka7/ava/src/config"
)

// search for songs in the server
func SearchServer(c *config.Connection, term ...string) (result SearchResult, err error) {
	err = validFilter(term...)
	config.Log(err)
	if err != nil {
		return nil, err
	}
	filter := term[0]
	result = make(map[string]Files)
	query, err := c.Client.Search(term...)
	config.Log(err)
	if err != nil {
		return nil, err
	}
	if len(query) >= 100 {
		query = query[:100]
	}
	for _, song := range query {

		index := song[filter]

		if _, duplicate := result[index]; !duplicate {
			result[index] = make(Files, 0)
		}

		result[index] = append(result[index], newFile(song))
	}
	return
}

// searchs for the term in server and adds them to current queue
func SearchAdd(c *config.Connection, term ...string) error {
	err := validFilter(term...)
	if err != nil {
		return err
	}
	err = searchAdd(c, term...)
	return err
}

// clears the queue then adds the founded songs to queue and plays it
func SearchPlay(c *config.Connection, term ...string) error {
	err := validFilter(term...)
	if err != nil {
		return err
	}
	err = c.Client.Clear()
	if err != nil {
		return err
	}
	err = searchAdd(c, term...)
	if err != nil {
		return err
	}
	c.Client.Play(-1)
	return nil
}

func searchAdd(c *config.Connection, term ...string) error {
	err = c.Client.Command("searchadd %s %s", term[0], term[1]).OK()
	return err
}

func validFilter(filter ...string) error {
	v := map[string]bool{
		"file": true, "Artist": true, "Album": true,
		"Genre": true, "Date": true, "Title": true,
	}
	if len(filter) < 2 {
		return errors.New("you have to enter a term")
	}
	if _, ok := v[filter[0]]; !ok {
		return errors.New("valid filters are : File,Artist,Album,Genre,Date and Title")
	}
	if len(filter[1]) == 0 {
		return errors.New("You have to enter something to query about")
	}
	return nil
}
