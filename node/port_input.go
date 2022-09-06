package node

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/uuid"
)

type Details struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type InputPort struct {
	*buffer.Const
	uuid        uuid.Value
	direction   Direction
	connector   *Connector
	NodeDetails *Details `json:"nodeDetails"`
}

func NewInputPort(_type buffer.Type, nodeDetails *Details) *InputPort {
	return &InputPort{buffer.NewConst(_type), uuid.New(), DirectionInput, nil, nodeDetails}
}

func (p *InputPort) UUID() uuid.Value {
	return p.uuid
}

func (p *InputPort) GetNodeDetails() *Details {
	return p.NodeDetails
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
