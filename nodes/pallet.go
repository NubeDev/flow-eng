package nodes

import (
	"github.com/NubeDev/flow-eng/node"
)

type PalletInputs struct {
	Name         string      `json:"name"`
	ValueType    string      `json:"valueType"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
}

type PalletOutputs struct {
	Name      string `json:"name"`
	ValueType string `json:"valueType"`
}

type Info struct {
	Icon string `json:"icon,omitempty"`
}

type PalletNode struct {
	Type          string           `json:"type"`
	Category      string           `json:"category"`
	IsParent      bool             `json:"isParent"`
	AllowSettings bool             `json:"allowSettings"`
	PalletInputs  []*PalletInputs  `json:"inputs"`
	PalletOutputs []*PalletOutputs `json:"outputs"`
	Info          Info             `json:"info"`
	Metadata      *node.Metadata   `json:"metadata,omitempty"`
	PayloadType   string           `json:"payloadType"`
	AllowPayload  bool             `json:"allowPayload"`
}

func convertOutputs(node *node.Spec) []*PalletOutputs {
	var all []*PalletOutputs
	for _, output := range node.GetOutputs() {
		one := &PalletOutputs{}
		one.Name = string(output.Name)
		one.ValueType = string(output.DataType)
		all = append(all, one)
	}
	return all
}
func convertInputs(node *node.Spec) []*PalletInputs {
	var all []*PalletInputs
	for _, input := range node.GetInputs() {
		one := &PalletInputs{}
		one.Name = string(input.Name)
		one.ValueType = string(input.DataType)
		all = append(all, one)
	}
	return all
}

func convertInfo(nodeInfo node.Info) Info {
	return Info{
		Icon: nodeInfo.Icon,
	}
}

func EncodePallet() ([]*PalletNode, error) {
	var all []*PalletNode
	for _, spec := range All() {
		one := &PalletNode{}
		nodeType, err := setType(spec)
		if err != nil {
			return nil, err
		}
		if spec.GetSchema() != nil {
			one.AllowSettings = true
		}
		if spec.AllowSettings {
			one.AllowSettings = true
		}
		one.Type = nodeType
		one.Category = spec.Info.Category
		one.IsParent = spec.IsParent
		one.Metadata = spec.GetMetadata()
		one.PalletInputs = convertInputs(spec)
		one.PalletOutputs = convertOutputs(spec)
		one.Info = convertInfo(spec.GetInfo())
		one.AllowPayload = spec.GetAllowPayload()
		one.PayloadType = string(spec.GetPayloadType())
		all = append(all, one)
	}
	return all, nil
}
