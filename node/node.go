package node

type Node interface {
	Process() // runs the logic of the node
	Cleanup()
	GetInfo() Info
	GetID() string       // node_abc123
	GetName() string     // AND, OR
	GetNodeName() string // my-node
	GetNodeValues() []*PortValues
	GetInputs() []*Input
	GetInput(name InputName) *Input
	GetOutputs() []*Output
	GetOutput(name OutputName) *Output
	InputUpdated(name InputName) bool
	InputsLen() int
	OutputsLen() int
	ReadPinAsInt(name InputName) int
	ReadPinAsFloat(name InputName) float64
	ReadPin(InputName) interface{}
	ReadMultiple(count int) []interface{}
	WritePin(OutputName, interface{})
	OverrideInputValue(name InputName, value interface{}) error
	GetMetadata() *Metadata
	GetIsParent() bool
	GetParameters() *Parameters
	GetSubFlow() *SubFlow
	GetSubFlowNodes() []*Spec
	DeleteSubFlowNodes()
	SetMetadata(m *Metadata)
	GetSettings() []*Settings
	GetSetting(name SettingTitle) *Settings
	SetPropertiesValue(name SettingTitle, value interface{}) error
}

func New(id, name, nodeName string, meta *Metadata) *Spec {
	return &Spec{
		Inputs:  nil,
		Outputs: nil,
		Info: Info{
			NodeID:   id,
			Name:     name,
			NodeName: nodeName,
		},
		Metadata: meta,
	}
}

type Spec struct {
	Inputs     []*Input    `json:"inputs,omitempty"`
	Outputs    []*Output   `json:"outputs,omitempty"`
	Info       Info        `json:"info"`
	Settings   []*Settings `json:"settings,omitempty"`
	Metadata   *Metadata   `json:"metadata,omitempty"`
	Parameters *Parameters `json:"parameters,omitempty"`
	IsParent   bool        `json:"isParent,omitempty"`
	SubFlow    *SubFlow    `json:"subFlow,omitempty"`
}

func (n *Spec) GetInfo() Info {
	return n.Info
}

func (n *Spec) GetID() string {
	return n.Info.NodeID
}

func (n *Spec) GetName() string {
	return n.Info.Name
}

func (n *Spec) GetIsParent() bool {
	return n.IsParent
}

func (n *Spec) GetNodeName() string {
	return n.Info.NodeName
}

func (n *Spec) GetInputs() []*Input {
	return n.Inputs
}

type PortValues struct {
	Type  DataTypes   `json:"type"`
	Value interface{} `json:"value"`
}

func (n *Spec) GetNodeValues() []*PortValues {
	var out []*PortValues
	for _, input := range n.Inputs {
		input.GetValue()
		out = append(out, &PortValues{
			Type:  input.DataType,
			Value: input.GetValue(),
		})
	}
	return out
}

func (n *Spec) GetInput(name InputName) *Input {
	for _, input := range n.GetInputs() {
		if input.Name == name {
			return input
		}
	}
	return nil
}

func (n *Spec) GetOutputs() []*Output {
	return n.Outputs
}

func (n *Spec) GetOutput(name OutputName) *Output {
	for _, out := range n.GetOutputs() {
		if out.Name == name {
			return out
		}
	}
	return nil
}

func (n *Spec) GetParameters() *Parameters {
	return n.Parameters
}

func (n *Spec) InputsLen() int {
	return len(n.Inputs)
}

func (n *Spec) OutputsLen() int {
	return len(n.Outputs)
}

func (n *Spec) GetSubFlow() *SubFlow {
	return n.SubFlow
}

func (n *Spec) GetSubFlowNodes() []*Spec {
	return n.SubFlow.Nodes
}

func (n *Spec) DeleteSubFlowNodes() {
	n.SubFlow.Nodes = nil
}

func (n *Spec) GetMetadata() *Metadata {
	return n.Metadata
}

func (n *Spec) SetMetadata(m *Metadata) {
	n.Metadata = m
}

type Info struct {
	NodeID      string `json:"nodeID"`             // a123
	Name        string `json:"name"`               // add, or
	NodeName    string `json:"nodeName,omitempty"` // my-node-abc
	Category    string `json:"category,omitempty"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

type DataTypes string
type InputName string

type OutputName string
type ApplicationName string // bacnet, mqtt

const (
	TypeString DataTypes = "string"
	TypeInt    DataTypes = "int"
	TypeFloat  DataTypes = "number"
	TypeNumber DataTypes = "number"
)

const (
	InputNamePrefix  string = "in"
	OutputNamePrefix string = "out"
)

const (
	InputCount SettingTitle = "input count"
	Operation  SettingTitle = "input count"
)

const (
	SetPoint InputName = "set-point"
	DeadBand InputName = "dead-band"

	Comment  InputName = "comment"
	InNumber InputName = "number"
	InString InputName = "string"

	In   InputName = "in"
	In1  InputName = "in1"
	In2  InputName = "in2"
	In3  InputName = "in3"
	In4  InputName = "in4"
	In10 InputName = "in10"
	In11 InputName = "in11"
	In12 InputName = "in12"
	In13 InputName = "in13"
	In14 InputName = "in14"
	In15 InputName = "in15"
	In16 InputName = "in16"

	Input_ InputName = "input"
	InputA InputName = "a"
	InputB InputName = "b"
	InputC InputName = "c"
	InputD InputName = "d"

	Topic InputName = "topic"

	DelaySeconds InputName = "delay (s)"
	Selection    InputName = "select"

	From InputName = "from"
	To   InputName = "to"

	Name          InputName = "name"
	ObjectId      InputName = "object-id"
	ObjectType    InputName = "object-type"
	OverrideInput InputName = "override-value"

	RisingEdge  InputName = "rising-edge"
	FallingEdge InputName = "falling-edge"
)

const (
	Result OutputName = "result"

	ErrMsg OutputName = "error"
	Msg    OutputName = "message"

	Toggle OutputName = "toggle"
	Out    OutputName = "out"

	OutNot OutputName = "out not"

	Out1 OutputName = "out1"
	Out2 OutputName = "out2"
	Out3 OutputName = "out3"
	Out4 OutputName = "out4"

	Above OutputName = "above"
	Below OutputName = "below"

	GraterThan OutputName = "grater"
	LessThan   OutputName = "less"
	Equal      OutputName = "equal"
)

type InputConnection struct {
	NodeID        string      `json:"nodeID,omitempty"`
	NodePort      OutputName  `json:"nodePort,omitempty"`
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"currentValue,omitempty"`
	FallbackValue interface{} `json:"fallbackValue,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type OutputConnection struct {
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node connection and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"currentValue,omitempty"`
	FallbackValue interface{} `json:"fallbackValue,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type Metadata struct {
	PositionX string `json:"positionX"`
	PositionY string `json:"positionY"`
}

type SubFlow struct {
	ParentID string  `json:"parentID,omitempty"` // nodeID eg: bacnet-server node
	Nodes    []*Spec `json:"nodes,omitempty"`    // bacnet-point
}

type Application struct {
	Application ApplicationName `json:"application,omitempty"` // eg: bacnet-point belongs to bacnet-server
	IsChild     bool            `json:"isChild"`
}

type Parameters struct {
	Application  *Application `json:"application"`
	MaxNodeCount int          `json:"maxNodeCount,omitempty"` // eg: bacnet-server node can only be added once
}
