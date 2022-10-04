package transformations

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type Scale struct {
	*node.Spec
}

func NewScale(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, scaleNode, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	minIn := node.BuildInput(node.InMin, node.TypeFloat, nil, body.Inputs)
	maxIn := node.BuildInput(node.InMax, node.TypeFloat, nil, body.Inputs)
	minOut := node.BuildInput(node.OutMin, node.TypeFloat, nil, body.Inputs)
	maxOut := node.BuildInput(node.OutMax, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(in, minIn, minOut, maxIn, maxOut)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Scale{body}, nil
}

func (inst *Scale) Process() {
	in := inst.ReadPinAsFloat(node.In)
	minIn := inst.ReadPinAsFloat(node.InMin)
	maxIn := inst.ReadPinAsFloat(node.InMax)
	minOut := inst.ReadPinAsFloat(node.OutMin)
	maxOut := inst.ReadPinAsFloat(node.OutMax)
	inst.WritePin(node.Out, float.Scale(in, minIn, maxIn, minOut, maxOut))
}

func (inst *Scale) Cleanup() {}
