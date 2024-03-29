package math

import (
	"github.com/NubeDev/flow-eng/node"
	"math"
)

type Absolute struct {
	*node.Spec
}

func NewAbsolute(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, absolute, category)
	in1 := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in1)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Absolute{body}, nil
}

func (inst *Absolute) Process() {
	in, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		output := math.Abs(in)
		inst.WritePinFloat(node.Out, output)
	}
}
