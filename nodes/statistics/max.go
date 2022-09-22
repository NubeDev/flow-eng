package statistics

import (
	"github.com/NubeDev/flow-eng/helpers/integer"
	"github.com/NubeDev/flow-eng/node"
)

type Max struct {
	*node.Spec
}

func NewMax(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, max, category)
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
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, count, integer.NonNil(buildCount.Min), integer.NonNil(buildCount.Max), body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Result, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return &Max{body}, nil
}

func (inst *Max) Process() {
	Process(inst)
}

func (inst *Max) Cleanup() {}