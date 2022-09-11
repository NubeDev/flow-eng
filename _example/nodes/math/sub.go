package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Sub struct {
	*node.BaseNode
}

func NewSub(body *node.BaseNode) (node.Node, error) {
	body = node.Defaults(body, add, category)
	buildCount, setting, count, err := inputsCount(body)
	if err != nil {
		return nil, err
	}
	settings, err := node.BuildSettings(setting)
	if err != nil {
		return nil, err
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.In, node.TypeFloat, nil, count, buildCount.Min, buildCount.Max, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return &Sub{body}, nil
}

func (inst *Sub) Process() {
	Process(inst)
}

func (inst *Sub) Cleanup() {}
