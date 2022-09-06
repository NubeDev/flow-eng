package node

import (
	"github.com/NubeDev/flow-eng/buffer/adapter"
)

type NodeProcess interface {
	Process()
	Cleanup()
	GetName() string
	GetID() string
	GetType() string
	GetInfo() Info
	GetInputs() []*TypeInput
	GetOutputs() []*TypeOutput
}

type Node struct {
	InputList  []*TypeInput  `json:"inputs"`
	OutputList []*TypeOutput `json:"outputs"`
	Info       Info          `json:"info"`
}

func (n *Node) GetType() string {
	return ""
}

func (n *Node) GetInfo() Info {
	return n.Info
}

func (n *Node) GetID() string {
	return n.Info.NodeID
}

func (n *Node) GetName() string {
	return n.Info.Name
}

func (n *Node) GetInputs() []*TypeInput {
	return n.InputList
}

func (n *Node) GetOutputs() []*TypeOutput {
	return n.OutputList
}

type Info struct {
	NodeID      string `json:"nodeID"` // abc
	Name        string `json:"name"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type DataTypes string
type PortName string

const (
	TypeAny     DataTypes = "any"
	TypeString  DataTypes = "string"
	TypeInt8    DataTypes = "int8"
	TypeFloat64 DataTypes = "float64"
)

const (
	In1  PortName = "in1"
	In2  PortName = "in2"
	Out1 PortName = "out1"
)

type Connection struct {
	NodeID   string `json:"nodeID"`
	NodePort string `json:"nodePortName"`
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
	ValueFloat64 *adapter.Float64 `json:"-"`
	ValueString  *adapter.String  `json:"-"`
}

type TypeOutput struct {
	*PortCommonOut
	*OutputPort
	ValueFloat64 *adapter.Float64 `json:"-"`
	ValueString  *adapter.String  `json:"-"`
}
