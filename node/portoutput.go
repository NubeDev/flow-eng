package node

import (
	"errors"
	"github.com/NubeDev/flow-eng/uuid"
)

var ErrTypesMismatch = errors.New("provided buffers types are different")

type OutputPort struct {
	Name        PortName            `json:"name"` // out1
	DataType    DataTypes           `json:"type"` // int8
	Connections []*OutputConnection `json:"connections"`
	Value       interface{}
	uuid        uuid.Value
	direction   Direction
	connectors  []*Connector
}

func newOutputPort(body *OutputPort) *OutputPort {
	return &OutputPort{
		body.Name,
		body.DataType,
		body.Connections,
		nil,
		uuid.New(),
		DirectionOutput,
		make([]*Connector, 0, 1)}
}

func (p *OutputPort) Write(value interface{}) {
	p.Value = value
	for i := 0; i < len(p.connectors); i++ {
		conn := p.connectors[i]
		conn.Notify()
	}
}

func (p *OutputPort) UUID() uuid.Value {
	return p.uuid
}

func (p *OutputPort) Direction() Direction {
	return p.direction
}

func (p *OutputPort) Connectors() []*Connector {
	return p.connectors
}

func (p *OutputPort) Copy(other *InputPort) error {
	if p.DataType != other.DataType {
		return ErrTypesMismatch
	}
	other.Value = p.Value
	return nil
}

func (p *OutputPort) Connect(inputs ...*InputPort) {
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		if err := p.connectInput(input); err != nil {
			panic(err)
		}
	}
}

func (p *OutputPort) connectInput(input *InputPort) error {
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
