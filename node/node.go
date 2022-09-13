package node

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
	GetSettings() []*Settings
	SetPropValue(name Title, value interface{}) error
	OverrideInputValue(name PortName, value interface{}) error
	InputsLen() int
	OutputsLen() int
	ReadMultipleNums(count int) []float64
	ReadMultiple(count int) []*Input
	ReadPinsNum(...PortName) []*RedMultiplePins
	ReadPinNum(PortName) (*float64, float64, bool)
	ReadPin(PortName) (*string, string)
	WritePin(PortName, interface{})
	WritePinNum(PortName, float64)
	SetMetadata(m *Metadata)
	GetMetadata() *Metadata
}

type BaseNode struct {
	Inputs   []*Input    `json:"inputs,omitempty"`
	Outputs  []*Output   `json:"outputs,omitempty"`
	Info     Info        `json:"info"`
	Settings []*Settings `json:"settings,omitempty"`
	Metadata *Metadata   `json:"metadata,omitempty"`
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

func (n *BaseNode) InputsLen() int {
	return len(n.Inputs)
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

func (n *BaseNode) OutputsLen() int {
	return len(n.Outputs)
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

func (n *BaseNode) WritePin(name PortName, value interface{}) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(value)
	}
}

func (n *BaseNode) WritePinNum(name PortName, value float64) {
	n.WritePin(name, value)
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
	In   PortName = "in"
	In1  PortName = "in1"
	In2  PortName = "in2"
	In3  PortName = "in3"
	In4  PortName = "in4"
	Out  PortName = "out"
	Out1 PortName = "out1"
	Out2 PortName = "out2"
	Out3 PortName = "out3"
	Out4 PortName = "out4"
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
	FallbackValue interface{} `json:"fallbackValue,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type Metadata struct {
	PositionX string `json:"positionX"`
	PositionY string `json:"positionY"`
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
