package node

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/uuid"
)

type InputPort struct {
	*buffer.Const
	uuid      uuid.Value
	direction Direction
	connector *Connector
}

func NewInputPort(_type buffer.Type) *InputPort {
	return &InputPort{buffer.NewConst(_type), uuid.New(), DirectionInput, nil}
}

func (p *InputPort) UUID() uuid.Value {
	return p.uuid
}

func (p *InputPort) Direction() Direction {
	return p.direction
}

func (p *InputPort) Connectors() []*Connector {
	if p.connector == nil {
		return []*Connector{}
	}
	return []*Connector{p.connector}
}
