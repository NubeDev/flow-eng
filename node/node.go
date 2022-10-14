package node

import (
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/schemas"
	"time"
)

type Node interface {
	Process() // runs the bool of the node
	Cleanup()
	Loop() (count uint64, firstLoop bool)
	AddDB(d db.DB)
	GetDB() db.DB
	SetSchema(schema *schemas.Schema)
	GetSchema() *schemas.Schema
	GetInfo() Info
	GetID() string       // node_abc123
	GetName() string     // AND, OR
	GetNodeName() string // my-node
	GetNodeValues() []*PortValues
	GetInputs() []*Input
	GetInput(name InputName) *Input
	GetOutputs() []*Output
	GetOutput(name OutputName) *Output
	InputUpdated(name InputName) (updated bool, boolCOV bool)
	InputHasConnection(name InputName) bool
	InputsLen() int
	OutputsLen() int
	ReadPin(InputName) interface{}
	ReadPinAsString(name InputName) (value string, null bool)
	ReadPinAsInt(name InputName) (value int, null bool)
	ReadPinAsBool(name InputName) (value bool, null bool)
	ReadPinAsFloat(name InputName) (value float64, null bool)
	ReadPinAsDuration(name InputName) (value time.Duration, null bool)
	ReadMultiple(count int) []interface{}
	ReadMultipleFloatPointer(count int) []*float64
	ReadMultipleFloat(count int) []float64
	WritePin(OutputName, interface{})
	WritePinFloat(name OutputName, value float64, precision ...int)
	WritePinBool(OutputName, bool)
	WritePinFalse(name OutputName)
	WritePinTrue(name OutputName)
	WritePinNull(OutputName)
	OverrideInputValue(name InputName, value interface{}) error
	GetMetadata() *Metadata
	GetIsParent() bool
	GetParameters() *Parameters
	GetSubFlow() *SubFlow
	GetSubFlowNodes() []*Spec
	DeleteSubFlowNodes()
	SetMetadata(m *Metadata)
	GetSettings() map[string]interface{}
	NodeValues() *Values
}

func New(id, name, nodeName string, meta *Metadata, settings map[string]interface{}) *Spec {
	n := &Spec{
		Inputs:  nil,
		Outputs: nil,
		Info: Info{
			NodeID:   id,
			Name:     name,
			NodeName: nodeName,
		},
		Metadata: meta,
		Settings: settings,
	}
	return n
}

type Spec struct {
	Inputs        []*Input               `json:"inputs,omitempty"`
	Outputs       []*Output              `json:"outputs,omitempty"`
	Info          Info                   `json:"info"`
	Settings      map[string]interface{} `json:"settings"`
	AllowSettings bool                   `json:"allowSettings"`
	Metadata      *Metadata              `json:"metadata,omitempty"`
	Parameters    *Parameters            `json:"parameters,omitempty"`
	IsParent      bool                   `json:"isParent,omitempty"`
	SubFlow       *SubFlow               `json:"subFlow,omitempty"`
	//OnStart       bool                   `json:"-"` // used for see if it's the first loop of the runner, if false it's the first run
	loopCount uint64
	schema    *schemas.Schema
	db        db.DB
}

func (n *Spec) Cleanup() {}

func (n *Spec) AddDB(d db.DB) {
	n.db = d
}

func (n *Spec) GetDB() db.DB {
	return n.db
}

// Loop will give you the loop count and a flag if it's the first loop
func (n *Spec) Loop() (loopCount uint64, firstLoop bool) {
	if n.loopCount == 0 {
		firstLoop = true
	} else {
		firstLoop = false
	}
	n.loopCount++
	return n.loopCount, firstLoop
}

func (n *Spec) GetSchema() *schemas.Schema {
	return n.schema
}

func (n *Spec) SetSchema(schema *schemas.Schema) {
	n.schema = schema
	n.AllowSettings = true
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

type InputConnection struct {
	NodeID        string      `json:"nodeID,omitempty"`
	NodePort      OutputName  `json:"nodePort,omitempty"`
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node link and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"currentValue,omitempty"`
	FallbackValue interface{} `json:"fallbackValue,omitempty"`
	Disable       *bool       `json:"disable,omitempty"`
}

type OutputConnection struct {
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node link and writes the value direct (or can be used to override a value)
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
	Application names.ApplicationName `json:"application,omitempty"` // eg: bacnet-point belongs to bacnet-server
	IsChild     bool                  `json:"isChild"`
}

type Parameters struct {
	Application  *Application `json:"application"`
	MaxNodeCount int          `json:"maxNodeCount,omitempty"` // eg: bacnet-server node can only be added once
}
