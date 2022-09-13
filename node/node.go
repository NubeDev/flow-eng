package node

type Node interface {
	Process() // runs the logic of the node
	Cleanup()
	GetInfo() Info
	GetID() string       // node_abc123
	GetName() string     // AND, OR
	GetNodeName() string // my-node
	GetInputs() []*Input
	GetInput(name InputName) *Input
	GetOutputs() []*Output
	GetOutput(name OutputName) *Output
	InputsLen() int
	OutputsLen() int
	ReadPin(InputName) interface{}
	ReadMultiple(count int) []interface{}
	WritePin(OutputName, interface{})
	OverrideInputValue(name InputName, value interface{}) error
	GetMetadata() *Metadata
	SetMetadata(m *Metadata)
	GetSettings() []*Settings
	SetPropertiesValue(name Title, value interface{}) error
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

func (n *BaseNode) GetName() string {
	return n.Info.Name
}

func (n *BaseNode) GetNodeName() string {
	return n.Info.NodeName
}

func (n *BaseNode) GetInputs() []*Input {
	return n.Inputs
}

func (n *BaseNode) GetInput(name InputName) *Input {
	for _, input := range n.GetInputs() {
		if input.Name == name {
			return input
		}
	}
	return nil
}

func (n *BaseNode) GetOutputs() []*Output {
	return n.Outputs
}

func (n *BaseNode) GetOutput(name OutputName) *Output {
	for _, out := range n.GetOutputs() {
		if out.Name == name {
			return out
		}
	}
	return nil
}

func (n *BaseNode) InputsLen() int {
	return len(n.Inputs)
}

func (n *BaseNode) OutputsLen() int {
	return len(n.Outputs)
}

func (n *BaseNode) GetMetadata() *Metadata {
	return n.Metadata
}

func (n *BaseNode) SetMetadata(m *Metadata) {
	n.Metadata = m
}

type Info struct {
	NodeID      string `json:"node_id"`             // a123
	Name        string `json:"name"`                // add, or
	NodeName    string `json:"node_name,omitempty"` // my-node-abc
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

type DataTypes string
type InputName string
type OutputName string

const (
	TypeString DataTypes = "string"
	TypeInt    DataTypes = "int"
	TypeFloat  DataTypes = "float"
)

const (
	InputNamePrefix  string = "in"
	OutputNamePrefix string = "out"
)

const (
	In1 InputName = "in1"
	In2 InputName = "in2"
	In3 InputName = "in3"
	In4 InputName = "in4"
)

const (
	Out1 OutputName = "out1"
	Out2 OutputName = "out2"
	Out3 OutputName = "out3"
	Out4 OutputName = "out4"
)

type InputConnection struct {
	NodeID        string      `json:"node_id,omitempty"`
	NodePort      OutputName  `json:"node_port_name,omitempty"`
	OverrideValue interface{} `json:"override_value,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"current_value,omitempty"`
	FallbackValue interface{} `json:"fallback_value,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type OutputConnection struct {
	NodeID        string      `json:"node_id,omitempty"`
	NodePort      InputName   `json:"node_port_name,omitempty"`
	OverrideValue interface{} `json:"override_value,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"current_value,omitempty"`
	FallbackValue interface{} `json:"fallback_value,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type Metadata struct {
	PositionX string `json:"position_x"`
	PositionY string `json:"position_y"`
}
