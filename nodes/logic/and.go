package logic

import (
	"github.com/NubeDev/flow-eng/node"
)

type And struct {
	*node.Spec
}

func NewAnd(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, and, category)
	buildCount, setting, value, err := node.NewSetting(body, &node.SettingOptions{Type: node.Number, Title: node.InputCount, Min: 2, Max: 20})
	if err != nil {
		return nil, err
	}
	settings, err := node.BuildSettings(setting)
	if err != nil {
		return nil, err
	}
	count, ok := value.(int)
	if !ok {
		count = 2
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, count, buildCount.Min, buildCount.Max, body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Result, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return &And{body}, nil
}

func (inst *And) Process() {
	Process(inst)
}

func (inst *And) Cleanup() {}