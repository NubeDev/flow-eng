package statistics

import (
	"github.com/NubeDev/flow-eng/node"
)

type Max struct {
	*node.Spec
}

func NewMax(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, max, category)
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 2, 2, 20, body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Max{body}, nil
}

func (inst *Max) Process() {
	Process(inst)
}
