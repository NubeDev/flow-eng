package node

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/str"
)

type Node interface {
	Process() // runs the logic of the node
	Cleanup()
	GetID() string       // node_abc123
	GetName() string     // AND, OR
	GetNodeName() string // my-node
	GetInfo() Info
	GetInputs() []*Input
	GetInput(name PortName) *Input
	GetOutputs() []*Output
	GetOutput(name PortName) *Output
	OverrideInputValue(name PortName, value interface{}) error
	ReadPinsNum(...PortName) []*RedMultiplePins
	ReadPinNum(PortName) (*float64, float64, bool)
	ReadPin(PortName) (*string, string)
	WritePin(PortName, *string)
	WritePinNum(PortName, float64)
	SetMetadata(m *Metadata)
	GetMetadata() *Metadata
}

type BaseNode struct {
	Inputs   []*Input    `json:"inputs,omitempty"`
	Outputs  []*Output   `json:"outputs,omitempty"`
	Info     Info        `json:"info"`
	Settings []*Settings `json:"settings,omitempty"`
	Metadata *Metadata   `json:"metadata"`
}

func (n *BaseNode) GetInfo() Info {
	return n.Info
}

func (n *BaseNode) GetID() string {
	return n.Info.NodeID
}

func (n *BaseNode) GetMetadata() *Metadata {
	return n.Metadata
}

func (n *BaseNode) SetMetadata(m *Metadata) {
	n.Metadata = m
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

func (n *BaseNode) GetInput(name PortName) *Input {
	for _, input := range n.GetInputs() {
		if input.Name == name {
			return input
		}
	}
	return nil
}

func (n *BaseNode) GetOutput(name PortName) *Output {
	for _, out := range n.GetOutputs() {
		if out.Name == name {
			return out
		}
	}
	return nil
}

func (n *BaseNode) GetOutputs() []*Output {
	return n.Outputs
}

type RedMultiplePins struct {
	Value *float64
	Real  float64
	Found bool
}

func (n *BaseNode) ReadPinsNum(name ...PortName) []*RedMultiplePins {
	var out []*RedMultiplePins
	var resp *RedMultiplePins
	for _, portName := range name {
		v, r, f := n.ReadPinNum(portName)
		resp.Value = v
		resp.Real = r
		resp.Found = f
		out = append(out, resp)
	}
	return out
}

func (n *BaseNode) OverrideInputValue(name PortName, value interface{}) error {
	in := n.GetInput(name)
	if in == nil {
		return errors.New(fmt.Sprintf("failed to find port%s", name))
	}
	if in.Connection != nil {
		in.Connection.OverrideValue = value
	} else {
		return errors.New(fmt.Sprintf("this node has no inputs"))
	}
	return nil

}

func (n *BaseNode) ReadPinNum(name PortName) (value *float64, real float64, hasValue bool) {
	pinValPointer, _ := n.ReadPin(name)
	valPointer, val, err := float.StringFloatErr(pinValPointer)
	if err != nil {
		return nil, 0, hasValue
	}
	return valPointer, val, float.NotNil(valPointer)
}

func (n *BaseNode) ReadPin(name PortName) (*string, string) {
	for _, out := range n.GetInputs() {
		if name == out.Name {
			if out.Connection.OverrideValue != nil { // this would be that the user wrote a value to the input directly
				toStr := fmt.Sprintf("%v", out.Connection.OverrideValue)
				return str.New(toStr), str.NonNil(str.New(toStr))
			}
			if out.Connection.FallbackValue != nil { // this would be that the user wrote a value to the input directly
				toStr := fmt.Sprintf("%v", out.Connection.FallbackValue)
				return str.New(toStr), str.NonNil(str.New(toStr))
			}
			val := out.Value.Get()
			return val, str.NonNil(val)
		}
	}
	return nil, ""
}

func (n *BaseNode) WritePin(name PortName, value *string) {
	for _, out := range n.GetOutputs() {
		if name == out.Name {
			out.Value.Set(value)
		}
	}
}

func (n *BaseNode) WritePinNum(name PortName, value float64) {
	n.WritePin(name, float.ToStrPtr(value))
}

type Info struct {
	NodeID      string `json:"nodeID"`             // a123
	Name        string `json:"name"`               // add, or
	NodeName    string `json:"nodeName,omitempty"` // my-node-abc
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

type DataTypes string
type PortName string

const (
	TypeString DataTypes = "string"
	TypeInt    DataTypes = "int"
	TypeFloat  DataTypes = "float"
)

const (
	Topic PortName = "topic"
	In1   PortName = "in1"
	In2   PortName = "in2"
	In3   PortName = "in3"
	In4   PortName = "in4"
	Out1  PortName = "out1"
	Out2  PortName = "out2"
	Out3  PortName = "out3"
	Out4  PortName = "out4"
)

type InputConnection struct {
	NodeID        string      `json:"nodeID,omitempty"`
	NodePort      PortName    `json:"nodePortName,omitempty"`
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"currentValue,omitempty"`
	FallbackValue interface{} `json:"fallbackValue,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type OutputConnection struct {
	NodeID        string      `json:"nodeID,omitempty"`
	NodePort      PortName    `json:"nodePortName,omitempty"`
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"currentValue,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type Metadata struct {
	PositionX string `json:"positionX"`
	PositionY string `json:"positionY"`
}

type Input struct {
	*InputPort
	Value *adapter.String `json:"-"`
}

type Output struct {
	*OutputPort
	Value *adapter.String `json:"-"`
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
