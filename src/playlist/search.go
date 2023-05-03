package playlist

import (
	"errors"

	"github.com/smoka7/ava/src/config"
)

var tags = map[string]struct{}{
	"file": {}, "Artist": {}, "Album": {},
	"Genre": {}, "Date": {}, "Title": {},
}

// search for songs in the server
func SearchServer(c config.Connection, term ...string) (SearchResult, error) {
	err = validFilter(term...)
	config.Log(err)
	if err != nil {
		return nil, err
	}
	// search the server
	query, err := c.Search(term...)
	config.Log(err)
	if err != nil {
		return nil, err
	}

	// only show Limited amount of results
	if len(query) >= ClientQueueLimit {
		query = query[:ClientQueueLimit]
	}

	// group the results by their searched Tag
	tag := term[0]
	result := make(map[string]Files, 0)
	for _, song := range query {

		index := song[tag]

		if _, duplicate := result[index]; !duplicate {
			result[index] = make(Files, 0)
		}

		result[index] = append(result[index], newFile(song))
	}
	return result, nil
}

// searchs for the term in server and adds them to current queue
func SearchAdd(c config.Connection, term ...string) error {
	err := validFilter(term...)
	if err != nil {
		return err
	}
	err = searchAdd(c, term...)
	return err
}

// clears the queue then adds the founded songs to queue and plays it
func SearchPlay(c config.Connection, term ...string) error {
	err := validFilter(term...)
	if err != nil {
		return err
	}
	err = c.Clear()
	if err != nil {
		return err
	}
	err = searchAdd(c, term...)
	if err != nil {
		return err
	}
	c.Play(-1)
	return nil
}

func searchAdd(c config.Connection, term ...string) error {
	err = c.Command("searchadd %s %s", term[0], term[1]).OK()
	return err
}

func validFilter(filters ...string) error {
	if len(filters) < 2 {
		return errors.New("you have to enter a term")
	}

	if len(filters)%2 != 0 {
		return errors.New("you have to write a value for every Tag")
	}

	for i := 0; i < len(filters); i += 2 {
		if _, ok := tags[filters[i]]; !ok {
			return errors.New("valid Tags are : File,Artist,Album,Genre,Date and Title")
		}

		if len(filters[i+1]) == 0 {
			return errors.New("You have to fill every Tag")
		}
	}
	return nil
}
