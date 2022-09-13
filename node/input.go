package node

import (
	"github.com/NubeDev/flow-eng/uuid"
)

type Input struct {
	Name       InputName        `json:"name"` // in1
	DataType   DataTypes        `json:"type"` // int
	Connection *InputConnection `json:"connection,omitempty"`
	value      interface{}
	uuid       uuid.Value
	direction  Direction
	connector  *Connector
}

func newInput(body *Input) *Input {
	return &Input{
		body.Name,
		body.DataType,
		body.Connection,
		nil,
		uuid.New(),
		DirectionInput,
		nil,
	}
}

func (p *Input) UUID() uuid.Value {
	return p.uuid
}

func (p *Input) Direction() Direction {
	return p.direction
}

func (p *Input) Connectors() []*Connector {
	if p.connector == nil {
		return []*Connector{}
	}
	return []*Connector{p.connector}
}

func (p *Input) GetValue() interface{} {
	return p.value
}

func (p *Input) SetValue(value interface{}) {
	p.value = value
}
