package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Divide struct {
	*node.Spec
}

func NewDivide(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, divide, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs, false, false)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in1, in2)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Divide{body}, nil
}

func (inst *Divide) Process() {
	in1, null1 := inst.ReadPinAsFloat(node.In1)
	in2, null2 := inst.ReadPinAsFloat(node.In2)
	if null1 || null2 || in2 == 0 {
		inst.WritePinNull(node.Out)
	} else {
		output := in1 / in2
		inst.WritePinFloat(node.Out, output)
	}

}
