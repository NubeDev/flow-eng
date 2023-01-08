package node

import (
	"github.com/NubeDev/flow-eng/schemas"
)

type Help struct {
	NodeName     string          `json:"name"`
	Help         string          `json:"help"`
	Info         Info            `json:"info,omitempty"`
	AllowPayload bool            `json:"allowPayload"`
	PayloadType  string          `json:"payloadType"`
	Inputs       []*Input        `json:"inputs,omitempty"`
	Outputs      []*Output       `json:"outputs,omitempty"`
	Settings     *schemas.Schema `json:"settings,omitempty"`
}

func (n *Spec) NodeHelp() *Help {
	inputs := n.GetInputs()
	for _, input := range inputs {
		input.Connection = nil
	}
	outputs := n.GetOutputs()
	for _, output := range outputs {
		output.Connections = nil
	}
	var out = &Help{
		NodeName:     n.GetName(),
		Info:         n.GetInfo(),
		Help:         n.GetHelp(),
		AllowPayload: n.GetAllowPayload(),
		PayloadType:  string(n.GetPayloadType()),
		Inputs:       inputs,
		Outputs:      outputs,
		Settings:     n.GetSchema(),
	}
	return out
}

const (
	InHelp            InputHelp = "interval"
	IntervalInputHelp InputHelp = "interval"
)

const (
	IntervalOutputHelp OutputHelp = "interval"
)
