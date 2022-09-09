package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers/float"
)

type Node interface {
	Process() // runs the logic of the node
	Cleanup()
	GetID() string       // node_abc123
	GetName() string     // AND, OR
	GetNodeName() string // my-node
	GetInfo() Info
	GetInputs() []*Input
	GetOutputs() []*Output
	//ReadPinsNum(...PortName) []*RedMultiplePins
	ReadPinNum(PortName) (interface{}, float64, bool)
	ReadPin(PortName) interface{}
	WritePin(PortName, interface{})
	WritePinNum(PortName, float64)
}

type BaseNode struct {
	Inputs   []*Input    `json:"inputs"`
	Outputs  []*Output   `json:"outputs"`
	Info     Info        `json:"info"`
	Settings []*Settings `json:"settings"`
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

type RedMultiplePins struct {
	Value *float64
	Real  float64
	Found bool
}

//func (n *BaseNode) ReadPinsNum(name ...PortName) []*RedMultiplePins {
//	var out []*RedMultiplePins
//	var resp *RedMultiplePins
//	for _, portName := range name {
//		v, r, f := n.ReadPinNum(portName)
//		resp.Value = v
//		resp.Real = r
//		resp.Found = f
//		out = append(out, resp)
//	}
//	return out
//}

func (n *BaseNode) ReadPinNum(name PortName) (value interface{}, real float64, ok bool) {
	pinValPointer := n.ReadPin(name)
	val, ok := pinValPointer.(float64)
	if !ok {
		return nil, 0, ok
	}
	return pinValPointer, val, ok
}

func (n *BaseNode) ReadPin(name PortName) interface{} {
	for _, out := range n.GetInputs() {
		if name == out.Name {
			if out.Connection.OverrideValue != nil { // this would be that the user wrote a value to the input directly
				fmt.Println(4444, out.Connection.OverrideValue)
				return out.Connection.OverrideValue
			}
			val := out.Value.Get()

			return val
		}

	}
	return nil
}

func (n *BaseNode) WritePin(name PortName, value interface{}) {
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
	NodeID        string      `json:"nodeID"`
	NodePort      PortName    `json:"nodePortName"`
	OverrideValue interface{} `json:"value,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"readValue,omitempty"`
	Disable       *bool       `json:"disable"`
}

type OutputConnection struct {
	NodeID        string      `json:"nodeID"`
	NodePort      PortName    `json:"nodePortName"`
	OverrideValue interface{} `json:"value,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"readValue,omitempty"`
	Disable       *bool       `json:"disable"`
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
