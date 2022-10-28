package node

import (
	"errors"
	"github.com/NubeDev/flow-eng/helpers/global"
	"github.com/NubeDev/flow-eng/helpers/uuid"
)

var (
	ErrNoInputData = errors.New("no input conversions was received")
)

type Connector struct {
	uuid    uuid.Value
	from    *Output
	to      *Input
	written bool
}

func NewConnector(from *Output, to *Input) *Connector {
	if from.DataType != to.DataType {
		panic(ErrIncompatiblePortsTypes)
	}
	return &Connector{uuid.New(), from, to, false}
}

func (connector *Connector) FromUUID() uuid.Value {
	return connector.from.UUID()
}

func (connector *Connector) ToUUID() uuid.Value {
	return connector.to.UUID()
}

func (connector *Connector) Trigger(debug *global.Debug) error {
	// exit if no new conversions was received
	if !connector.written {
		return ErrNoInputData
	}

	// move conversions to destination port
	err := connector.from.Copy(connector.to, debug)
	return err
}

func (connector *Connector) Notify() {
	connector.written = true
}

func (connector *Connector) Reset() {
	connector.written = false
}
