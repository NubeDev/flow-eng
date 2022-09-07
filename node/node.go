package node

import (
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/str"
	"strconv"
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

func (n *BaseNode) readPinNum(name PortName) (*float64, float64, bool) {
	pinValPointer, _ := n.readPin(name)
	valPointer, val, err := float.StringFloatErr(pinValPointer)
	if err != nil {
		return nil, 0, true
	}
	return valPointer, val, float.NotNil(valPointer)
}

func (n *BaseNode) readPin(name PortName) (*string, string) {
	for _, out := range n.GetInputs() {
		if name == out.Name {
			if !str.IsNil(out.Connection.Value) { // this would be that the user wrote a value to the input directly
				return out.Connection.Value, str.NonNil(out.Connection.Value)
			}
			val := out.Value.Get()
			return val, str.NonNil(val)
		}
	}
	return nil, ""
}

func (n *BaseNode) writePin(name PortName, value *string) bool {
	for _, out := range n.GetOutputs() {
		if name == out.Name {
			out.Value.Set(value)
			return true
		}
	}
	return false
}

func (n *BaseNode) writePinNum(name PortName, value float64) bool {
	ok := n.writePin(name, float.ToStrPtr(value))
	return ok
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
	In1  PortName = "in1"
	In2  PortName = "in2"
	In3  PortName = "in3"
	In4  PortName = "in4"
	Out1 PortName = "out1"
	Out2 PortName = "out2"
	Out3 PortName = "out3"
	Out4 PortName = "out4"
)

type InputConnection struct {
	NodeID   string   `json:"nodeID"`
	NodePort PortName `json:"nodePortName"`
	Value    *string  `json:"value,omitempty"` // used for when the user has no node connection and writes the value direct
}

type Connection struct {
	NodeID   string   `json:"nodeID"`
	NodePort PortName `json:"nodePortName"`
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

func toFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case bool:
		return float64(BoolToInt(t))
	case int:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case float64:
		return t
	case float32:
		return float64(t)
	case string:
		if i, err := strconv.ParseFloat(t, 64); err == nil {
			return i
		}
	}

	return 0
}
