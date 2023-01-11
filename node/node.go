package node

import (
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/store"
	"github.com/NubeDev/flow-eng/schemas"
	"time"
)

type CurrentState int

const (
	CREATED CurrentState = iota
	STARTED CurrentState = iota
	STOPPED CurrentState = iota
)

type Node interface {
	Start()
	Process()
	Stop()

	GetCurrentState() CurrentState
	SetCurrentState(currentState CurrentState)

	ResetProcessed()
	SetProcessed()
	GetProcessed() bool

	Loop() (count uint64, firstLoop bool)
	AddDB(d db.DB)
	GetDB() db.DB
	AddStore(s *store.Store)
	GetStore() *store.Store
	SetSchema(schema *schemas.Schema)
	GetSchema() *schemas.Schema
	GetInfo() Info
	GetID() string       // node_abc123
	GetName() string     // AND, OR
	GetNodeName() string // my-node
	SetNodeName(string)  // my-node
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
	ReadPinOrSettings(name InputName) interface{}
	ReadPinAsString(name InputName) (value string, null bool)
	ReadPinAsInt(name InputName) (value int, null bool)
	ReadPinAsBool(name InputName) (value bool, null bool)
	ReadPinAsFloat(name InputName) (value float64, null bool)
	ReadInputPriority(name InputName) (value float64, null bool)
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
	GetParentId() string
	SetMetadata(m *Metadata)
	GetAllowSettings() bool
	SetAllowSettings()
	GetSettings() map[string]interface{}
	NodeValues() *Values
	GetStatus() *Status
	SetStatus(*Status)
	SetSubTitle(message string)
	SetStatusError(message string)
	SetStatusMessage(message string)
	SetNotifyMessage(message string)
	SetWaringMessage(message string)
	SetErrorIcon(icon string)
	SetNotifyIcon(icon string)
	SetWaringIcon(icon string)
	GetPayload() *Payload
	SetPayload(payload *Payload)
	ReadPayloadAsFloat() (value float64, null bool)
	GetPayloadNull() (value any, null bool)
	GetNode(uuid string) Node
	GetNodes() []Node
	AddNodes(f []Node)
	SetIcon(icon string)
	GetIcon() string
	SetHelp(string)
	GetHelp() string
	SetAllowPayload()
	GetAllowPayload() bool
	SetPayloadType(t DataTypes)
	GetPayloadType() DataTypes
	SetDynamicInputs()
	SetDynamicOutputs()
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
	Metadata      *Metadata              `json:"metadata"`
	IsParent      bool                   `json:"isParent"`
	ParentId      string                 `json:"parentId,omitempty"`
	Status        *Status                `json:"status,omitempty"`
	AllowPayload  bool                   `json:"allowPayload"` // allow payload write
	PayloadType   DataTypes              `json:"payloadType"`
	Payload       *Payload               `json:"payload,omitempty"`
	Help          string                 `json:"help"`
	loopCount     uint64
	schema        *schemas.Schema
	db            db.DB
	store         *store.Store
	nodes         []Node
	processed     bool
	currentState  CurrentState
}

func (n *Spec) Start() {}

func (n *Spec) Stop() {}

func (n *Spec) ResetProcessed() {
	n.processed = false
}

func (n *Spec) SetNodeName(name string) {
	n.Info.NodeName = name
}

func (n *Spec) SetProcessed() {
	n.processed = true
}

func (n *Spec) GetProcessed() bool {
	return n.processed
}

func (n *Spec) GetCurrentState() CurrentState {
	return n.currentState
}

func (n *Spec) SetCurrentState(currentState CurrentState) {
	n.currentState = currentState
}

func (n *Spec) SetIcon(icon string) {
	n.Info.Icon = icon
}

func (n *Spec) GetIcon() string {
	return n.Info.Icon
}

func (n *Spec) AddDB(d db.DB) {
	n.db = d
}

func (n *Spec) GetNode(uuid string) Node {
	for _, node := range n.nodes {
		if node.GetID() == uuid {
			return node
		}
	}
	return nil
}

func (n *Spec) GetNodes() []Node {
	return n.nodes
}

func (n *Spec) AddNodes(f []Node) {
	n.nodes = f
}

func (n *Spec) GetDB() db.DB {
	return n.db
}

func (n *Spec) AddStore(s *store.Store) {
	n.store = s
}

func (n *Spec) GetStore() *store.Store {
	return n.store
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

func (n *Spec) GetAllowSettings() bool {
	return n.AllowSettings
}

func (n *Spec) SetAllowSettings() {
	n.AllowSettings = true
}

func (n *Spec) SetSchema(schema *schemas.Schema) {
	n.schema = schema
	n.AllowSettings = true
}

func (n *Spec) GetInfo() Info {
	return n.Info
}

func (n *Spec) SetHelp(body string) {
	n.Help = body
}

func (n *Spec) GetHelp() string {
	return n.Help
}

func (n *Spec) GetID() string {
	return n.Info.NodeID
}

func (n *Spec) GetName() string {
	return n.Info.Name
}

func (n *Spec) SetAllowPayload() {
	n.AllowPayload = true
}

func (n *Spec) GetAllowPayload() bool {
	return n.AllowPayload
}

func (n *Spec) SetPayloadType(t DataTypes) {
	n.PayloadType = t
}

func (n *Spec) GetPayloadType() DataTypes {
	return n.PayloadType
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

func (n *Spec) InputsLen() int {
	return len(n.Inputs)
}

func (n *Spec) OutputsLen() int {
	return len(n.Outputs)
}

func (n *Spec) GetParentId() string {
	return n.ParentId
}

func (n *Spec) GetDisplay() string {
	return n.Info.Display
}

func (n *Spec) SetDisplay(body string) {
	n.Info.Display = body
}

type Info struct {
	NodeID      string `json:"nodeID"`             // a123
	Name        string `json:"name,omitempty"`     // add, or
	NodeName    string `json:"nodeName,omitempty"` // my-node-abc
	Category    string `json:"category,omitempty"`
	Type        string `json:"type,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Display     string `json:"display,omitempty"`
}

type InputConnection struct {
	NodeID        string      `json:"nodeID,omitempty"`
	NodePort      OutputName  `json:"nodePort,omitempty"`
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node link and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"currentValue,omitempty"`
	FallbackValue interface{} `json:"fallbackValue,omitempty"`
	Disable       bool        `json:"disable,omitempty"`
}

type OutputConnection struct {
	OverrideValue interface{} `json:"overrideValue,omitempty"` // used for when the user has no node link and writes the value direct (or can be used to override a value)
	CurrentValue  interface{} `json:"currentValue,omitempty"`
	FallbackValue interface{} `json:"fallbackValue,omitempty"`
	Disable       bool        `json:"disable,omitempty"`
}

type Metadata struct {
	PositionX      string `json:"positionX,omitempty"`
	PositionY      string `json:"positionY,omitempty"`
	DynamicInputs  bool   `json:"dynamicInputs"`
	DynamicOutputs bool   `json:"dynamicOutputs"`
}
