package node

import (
	"errors"
	"github.com/NubeDev/flow-eng/helpers/global"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/helpers/uuid"
)

var ErrTypesMismatch = errors.New("provided types are different")

type Output struct {
	Name         OutputName          `json:"name"` // out1
	DataType     DataTypes           `json:"type"` // int8
	Connections  []*OutputConnection `json:"connections,omitempty"`
	value        interface{}
	uuid         uuid.Value
	direction    Direction
	connectors   []*Connector
	Help         OutputHelp `json:"help"`
	FolderExport bool       `json:"folderExport"`
	HideOutput   bool       `json:"hideOutput"`
}

func newOutput(body *Output) *Output {
	return &Output{
		body.Name,
		body.DataType,
		body.Connections,
		nil,
		uuid.New(),
		DirectionOutput,
		make([]*Connector, 0, 1),
		"",
		false,
		false,
	}
}

func (p *Output) Write(value interface{}) {
	p.SetValue(value)
	for i := 0; i < len(p.connectors); i++ {
		conn := p.connectors[i]
		conn.Notify()
	}
}

func (p *Output) UUID() uuid.Value {
	return p.uuid
}

func (p *Output) Direction() Direction {
	return p.direction
}

func (p *Output) Connectors() []*Connector {
	return p.connectors
}

func (p *Output) GetValue() interface{} {
	return p.value
}

func (p *Output) SetValue(value interface{}) {
	p.value = value
}

func (p *Output) Copy(other *Input, debug *global.Debug) error {
	if !compatibleTypes(p.DataType, other.DataType) {
		return ErrTypesMismatch
	}
	if global.DebugConnections {
		debug.OutputValue = p.GetValue()
		pprint.PrintJSON(debug)
	}
	// pass the output value to the input
	other.SetValue(p.GetValue())
	return nil
}

func (p *Output) Connect(inputs ...*Input) {
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		if err := p.connectInput(input); err != nil {
			panic(err)
			// log.Error(err)
		}
	}
}

func (p *Output) connectInput(input *Input) error {
	if input.connector != nil {
		return ErrIncompatiblePorts
	}
	connector := NewConnector(p, input)
	// add connector to the destination port
	input.connector = connector
	// add connector to output port to enable notifiers
	p.connectors = append(p.connectors, connector)
	return nil
}
