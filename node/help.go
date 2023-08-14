package node

import (
	"fmt"
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

type inputHelpBuilder struct {
	Name      InputName
	Type      DataTypes
	Help      InputHelp
	Described string
}

func (n *Spec) buildHelpInput() []inputHelpBuilder {
	var out []inputHelpBuilder
	for _, input := range n.GetInputs() {
		i := inputHelpBuilder{
			Name:      input.Name,
			Type:      input.DataType,
			Help:      input.Help,
			Described: fmt.Sprintf("name: %s type: %s %s", input.Name, input.DataType, input.Help),
		}
		out = append(out, i)
	}
	return out
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
