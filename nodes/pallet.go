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

type PalletNode struct {
	Type          string           `json:"type"`
	Category      string           `json:"category"`
	PalletInputs  []*PalletInputs  `json:"inputs"`
	PalletOutputs []*PalletOutputs `json:"outputs"`
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

func EncodePalle2t() ([]*node.Spec, error) {
	var all []*node.Spec
	for _, spec := range All() {
		nodeType, err := setType(spec)
		if err != nil {
			return nil, err
		}
		spec.Info.Type = nodeType
		all = append(all, spec)

	}
	return all, nil
}

func EncodePallet() ([]*PalletNode, error) {
	var all []*PalletNode
	for _, spec := range All() {
		one := &PalletNode{}
		nodeType, err := setType(spec)
		if err != nil {
			return nil, err
		}
		one.Type = nodeType
		one.Category = spec.Info.Category
		one.PalletInputs = convertInputs(spec)
		one.PalletOutputs = convertOutputs(spec)
		all = append(all, one)

	}
	return all, nil
}
