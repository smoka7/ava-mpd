package config

import (
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
)

type (
	Outputs []Output
	Output  struct {
		ID      int
		Name    string
		Enabled bool
	}
)

func newOutput(value mpd.Attrs) Output {
	id, _ := strconv.Atoi(value["outputid"])
	return Output{
		ID:      id,
		Name:    value["outputname"],
		Enabled: value["outputenabled"] == "1",
	}
}

// returns the mpd outputs
func (c Connection) ListOutputs() Outputs {
	c.Connect()
	list, err := c.Client.ListOutputs()
	defer c.Close()

	Log(err)
	if err != nil {
		return nil
	}

	outputs := make(Outputs, len(list))
	for i, output := range list {
		outputs[i] = newOutput(output)
	}

	return outputs
}

func (c Connection) ToggleOutput(id int, enabled bool) {
	if enabled {
		err := c.Client.EnableOutput(id)
		Log(err)
		return
	}
	err := c.Client.DisableOutput(id)
	Log(err)
}
