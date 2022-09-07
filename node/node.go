package node

import (
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers/float"
)

type Node interface {
	Process()
	Cleanup()
	GetName() string // AND, OR
	GetNodeName() string
	GetInfo() Info
	GetID() string
	GetInputs() []*Input
	GetOutputs() []*Output
}

type BaseNode struct {
	Inputs  []*Input  `json:"inputs"`
	Outputs []*Output `json:"outputs"`
	Info    Info      `json:"info"`
}

func (n *BaseNode) GetInfo() Info {
	return n.Info
}

func (n *BaseNode) GetID() string {
	return n.Info.NodeID
}

func (n *BaseNode) GetName() string {
	return n.Info.Name
}

func (n *BaseNode) GetNodeName() string {
	return n.Info.NodeName
}

func (n *BaseNode) GetInputs() []*Input {
	return n.Inputs
}

func (n *BaseNode) GetOutputs() []*Output {
	return n.Outputs
}

func (n *BaseNode) readPinValue(name PortName) (*float64, float64) {
	for _, out := range n.GetInputs() {
		if name == out.Name {
			if !float.IsNil(out.Connection.Value) { // this would be that the user wrote a value to the input directly
				return out.Connection.Value, float.NonNil(out.Connection.Value)
			}
			val := out.ValueFloat64.Get()
			return val, float.NonNil(val)
		}
	}
	return nil, 0
}

func (n *BaseNode) writePinValue(name PortName, value *float64) bool {
	for _, out := range n.GetOutputs() {
		if name == out.Name {
			out.ValueFloat64.Set(value)
			return true
		}
	}
	return false
}

type Info struct {
	NodeID      string `json:"nodeID"`   // a123
	Name        string `json:"name"`     // add, or
	NodeName    string `json:"nodeName"` // my-node-abc
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

type InputConnection struct {
	NodeID   string   `json:"nodeID"`
	NodePort PortName `json:"nodePortName"`
	Value    *float64 `json:"value,omitempty"` // used for when the user has no node connection and writes the value direct
}

type Connection struct {
	NodeID   string   `json:"nodeID"`
	NodePort PortName `json:"nodePortName"`
}

type Input struct {
	// *PortCommon
	*InputPort
	ValueFloat64 *adapter.Float64 `json:"-"`
	ValueString  *adapter.String  `json:"-"`
}

type Output struct {
	*OutputPort
	ValueFloat64 *adapter.Float64 `json:"-"`
	ValueString  *adapter.String  `json:"-"`
}
