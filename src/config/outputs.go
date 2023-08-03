package config

import (
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
)

type (
	Outputs []Output
	Output  struct {
		Name    string
		ID      int
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

// returns the mpd outputs.
func (c Connection) GetOutputs() Outputs {
	c.Connect()
	list, err := c.ListOutputs()
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

func (c Connection) ToggleOutput(id int) {
	err := c.Command("toggleoutput %d", id).OK()
	Log(err)
}
