package point

import (
	"github.com/NubeDev/flow-eng/node"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pointNumber, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in1, in2)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &Number{body}, nil
}

func (inst *Number) Process() {
	in1, in1Null := inst.ReadPinAsFloat(node.In1)
	in2, in2Null := inst.ReadPinAsFloat(node.In2)
	if !in1Null {
		inst.WritePinFloat(node.Out, in1)
		return
	}
	if !in2Null {
		inst.WritePinFloat(node.Out, in2)
		return
	}
	inst.WritePinNull(node.Out)
}
