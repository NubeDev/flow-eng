package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Sub struct {
	*node.Spec
}

func NewSub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, add, category)
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
	return &Sub{body}, nil
}

func (inst *Sub) Process() {
	Process(inst)
}

func (inst *Sub) Cleanup() {}
