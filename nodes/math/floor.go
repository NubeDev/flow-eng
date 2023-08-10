package math

import (
	"math"

	"github.com/NubeDev/flow-eng/node"
)

type Floor struct {
	*node.Spec
}

func NewFloor(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, floor, Category)
	in1 := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in1)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Floor{body}, nil
}

func (inst *Floor) Process() {
	in, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		output := math.Floor(in)
		inst.WritePinFloat(node.Out, output)
	}
}
