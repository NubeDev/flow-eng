package transformations

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type Limit struct {
	*node.Spec
}

func NewLimit(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, limitNode, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	min := node.BuildInput(node.Min, node.TypeFloat, nil, body.Inputs)
	max := node.BuildInput(node.Max, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(in, min, max)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Limit{body}, nil
}

func (inst *Limit) Process() {
	in := inst.ReadPinAsFloat(node.In)
	min := inst.ReadPinAsFloat(node.Min)
	max := inst.ReadPinAsFloat(node.Max)
	inst.WritePin(node.Out, float.LimitToRange(in, min, max))
}

func (inst *Limit) Cleanup() {}
