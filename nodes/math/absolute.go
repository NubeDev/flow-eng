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
	in1 := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in1)

	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Absolute{body}, nil
}

func (inst *Absolute) Process() {
	in, null := inst.ReadPinAsFloat(node.In1)
	if null {
		inst.WritePinNull(node.Outp)
	} else {
		output := math.Abs(in)
		inst.WritePinFloat(node.Outp, output)
	}
}
