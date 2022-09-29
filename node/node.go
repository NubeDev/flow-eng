package node

import (
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/schemas"
)

type Node interface {
	Process() // runs the bool of the node
	Cleanup()
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
	InputUpdated(name InputName) bool
	InputsLen() int
	OutputsLen() int
	ReadPinAsString(name InputName) string
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
	GetSettings() map[string]interface{}
}

func New(id, name, nodeName string, meta *Metadata, settings map[string]interface{}) *Spec {
	return &Spec{
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
}

type Spec struct {
	Inputs        []*Input               `json:"inputs,omitempty"`
	Outputs       []*Output              `json:"outputs,omitempty"`
	Info          Info                   `json:"info"`
	Settings      map[string]interface{} `json:"settings,omitempty"`
	AllowSettings bool                   `json:"allowSettings"`
	Metadata      *Metadata              `json:"metadata,omitempty"`
	Parameters    *Parameters            `json:"parameters,omitempty"`
	IsParent      bool                   `json:"isParent,omitempty"`
	SubFlow       *SubFlow               `json:"subFlow,omitempty"`
	OnStart       bool                   `json:"-"` // used for see if it's the first loop of the runner, if false it's the first run
	schema        *schemas.Schema
	db            db.DB
}

func (n *Spec) AddDB(d db.DB) {
	n.db = d
}

func (n *Spec) GetDB() db.DB {
	return n.db
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
