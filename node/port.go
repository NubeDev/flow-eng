package node

import (
	"errors"
	"github.com/NubeDev/flow-eng/uuid"
)

var (
	ErrIncompatiblePorts = errors.New("incompatible ports provided")
)

type Direction int

const (
	DirectionInput  Direction = 0
	DirectionOutput Direction = 1
)

type Port interface {
	UUID() uuid.Value
	Direction() Direction
	Connectors() []*Connector
}
