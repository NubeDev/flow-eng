package node

import "github.com/NubeDev/flow-eng/buffer/adapter"

//type NodeProcess interface {
//	Process()
//	Cleanup()
//	GetName() string // AND, OR
//	GetInfo() Info
//	GetInputs() []*TypeInput
//	GetOutputs() []*TypeOutput
//}

type Node struct {
	Inputs  []*Input  `json:"inputs"`
	Outputs []*Output `json:"outputs"`
	Info    Info      `json:"info"`
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

func (n *Node) GetInputs() []*Input {
	return n.Inputs
}

func (n *Node) GetOutputs() []*Output {
	return n.Outputs
}

type Info struct {
	NodeID      string `json:"nodeID"` // abc
	Name        string `json:"name"`
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

type Input struct {
	//*PortCommon
	*InputPort
	ValueFloat64 *adapter.Float64 `json:"-"`
	ValueString  *adapter.String  `json:"-"`
}

type Output struct {
	*OutputPort
	ValueFloat64 *adapter.Float64 `json:"-"`
	ValueString  *adapter.String  `json:"-"`
}
