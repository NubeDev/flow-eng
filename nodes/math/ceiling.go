package math

import (
	"math"

	"github.com/NubeDev/flow-eng/node"
)

type Ceiling struct {
	*node.Spec
}

func NewCeiling(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, ceiling, Category)
	in1 := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in1)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Ceiling{body}, nil
}

func (inst *Ceiling) Process() {
	in, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		output := math.Ceil(in)
		inst.WritePinFloat(node.Out, output)
	}
}
