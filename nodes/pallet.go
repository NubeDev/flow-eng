package nodes

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
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
	one := &PalletOutputs{}
	for _, output := range node.GetOutputs() {
		one.Name = string(output.Name)
		one.ValueType = string(output.DataType)
		all = append(all, one)
	}
	return all
}
func convertInputs(node *node.Spec) []*PalletInputs {
	var all []*PalletInputs
	one := &PalletInputs{}
	for _, input := range node.GetInputs() {
		one.Name = string(input.Name)
		one.ValueType = string(input.DataType)
		all = append(all, one)
	}
	return all
}

func EncodePallet() ([]*PalletNode, error) {
	var all []*PalletNode
	one := &PalletNode{}
	for _, spec := range All() {
		pprint.Print(spec)
		nodeType, err := setType(spec)
		if err != nil {
			fmt.Println(err)
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
