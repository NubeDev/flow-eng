package numtransform

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type Scale struct {
	*node.Spec
}

func NewScale(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, scaleNode, category)
	in := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, false)
	minIn := node.BuildInput(node.InMin, node.TypeFloat, nil, body.Inputs, false)
	maxIn := node.BuildInput(node.InMax, node.TypeFloat, nil, body.Inputs, false)
	minOut := node.BuildInput(node.OutMin, node.TypeFloat, nil, body.Inputs, false)
	maxOut := node.BuildInput(node.OutMax, node.TypeFloat, nil, body.Inputs, false)
	inputs := node.BuildInputs(in, minIn, maxIn, minOut, maxOut)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Scale{body}, nil
}

func (inst *Scale) Process() {
	in, _ := inst.ReadPinAsFloat(node.Inp)
	minIn, _ := inst.ReadPinAsFloat(node.InMin)
	maxIn, _ := inst.ReadPinAsFloat(node.InMax)
	minOut, _ := inst.ReadPinAsFloat(node.OutMin)
	maxOut, _ := inst.ReadPinAsFloat(node.OutMax)
	inst.WritePinFloat(node.Out, float.Scale(in, minIn, maxIn, minOut, maxOut))
}
