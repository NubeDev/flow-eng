package statistics

import (
	"github.com/NubeDev/flow-eng/node"
)

type Min struct {
	*node.Spec
}

func NewMin(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, min, category)
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 2, 2, 20, body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Result, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Min{body}, nil
}

func (inst *Min) Process() {
	Process(inst)
}
