package pallet

import (
	"errors"
	"fmt"
	"strings"

	"github.com/NubeDev/flow-eng/node"
)

type PalletInputs struct {
	Name            string      `json:"name"`
	ValueType       string      `json:"valueType"`
	DefaultValue    interface{} `json:"defaultValue,omitempty"`
	FolderExport    bool        `json:"folderExport"`
	HideInput       bool        `json:"hideInput"`
	PreventOverride bool        `json:"preventOverride"`
}

type PalletOutputs struct {
	Name         string `json:"name"`
	ValueType    string `json:"valueType"`
	FolderExport bool   `json:"folderExport"`
	HideOutput   bool   `json:"hideOutput"`
}

type PalletNode struct {
	Type          string           `json:"type"`
	Category      string           `json:"category"`
	IsParent      bool             `json:"isParent"`
	AllowSettings bool             `json:"allowSettings"`
	PalletInputs  []*PalletInputs  `json:"inputs"`
	PalletOutputs []*PalletOutputs `json:"outputs"`
	Info          node.Info        `json:"info"`
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
		one.FolderExport = output.FolderExport
		one.HideOutput = output.HideOutput
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
		one.FolderExport = input.FolderExport
		one.HideInput = input.HideInput
		one.DefaultValue = input.Connection.DefaultValue
		one.PreventOverride = input.PreventOverride
		all = append(all, one)
	}
	return all
}

func convertInfo(nodeInfo node.Info) node.Info {
	return node.Info{
		Icon: nodeInfo.Icon,
	}
}

func EncodePallet() ([]*PalletNode, error) {
	var all []*PalletNode
	for _, spec := range All() {
		one := &PalletNode{}
		nodeType, err := SetType(spec)
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

func SetType(n *node.Spec) (string, error) {
	if n == nil {
		return "", errors.New("node info can not be empty")
	}
	if n.Info.Name == "" {
		return "", errors.New("node name can not be empty")
	}
	if n.Info.Category == "" {
		return "", errors.New("node category can not be empty")
	}
	return fmt.Sprintf("%s/%s", n.Info.Category, n.Info.Name), nil

}

func DecodeType(nodeType string) (category, name string, err error) {
	parts := strings.Split(nodeType, "/")
	if len(parts) > 1 {
		return parts[0], parts[1], nil
	}
	return "", "", errors.New("failed to get category and name from node-type")
}
