package math

import (
	"github.com/NubeDev/flow-eng/node"
	"math"
)

type Ceiling struct {
	*node.Spec
}

func NewCeiling(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, ceiling, category)
	in1 := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in1)

	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Ceiling{body}, nil
}

func (inst *Ceiling) Process() {
	in, null := inst.ReadPinAsFloat(node.In1)
	if null {
		inst.WritePinNull(node.Outp)
	} else {
		output := math.Ceil(in)
		inst.WritePinFloat(node.Outp, output)
	}
}
