package node

import "github.com/NubeDev/flow-eng/buffer/adapter"

type NodeInfo struct {
	Name        string `json:"-"`
	Type        string `json:"-"`
	Description string `json:"-"`
	Version     string `json:"-"`
}

type DataTypes string
type PortName string

const (
	TypeInt8 DataTypes = "int8"
)

const (
	In1  PortName = "in1"
	Out1 PortName = "out1"
)

type Connection struct {
	NodeID string `json:"nodeID"`
	Port   string `json:"port"`
}

type PortCommon struct {
	Name       PortName    `json:"name"` // in1
	Type       DataTypes   `json:"type"` // int8
	Connection *Connection `json:"connection"`
}

type PortCommonOut struct {
	Name        PortName      `json:"name"` // in1
	Type        DataTypes     `json:"type"` // int8
	Connections []*Connection `json:"connection"`
}

type TypeInput struct {
	*PortCommon
	*InputPort
	Value *adapter.Int8
}

type TypeOutput struct {
	*PortCommonOut
	*OutputPort
	Value *adapter.Int8
}

type Node interface {
	Process()
	Cleanup()
	Info() NodeInfo
	Inputs() []*TypeInput
	Outputs() []*TypeOutput
}
