package boolean

import (
	"github.com/NubeDev/flow-eng/node"
)

type Xor struct {
	*node.Spec
}

func NewXor(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, xor, category)
	in1 := node.BuildInput(node.In1, node.TypeBool, nil, body.Inputs, false, true)
	in2 := node.BuildInput(node.In2, node.TypeBool, nil, body.Inputs, false, true)
	inputs := node.BuildInputs(in1, in2)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Xor{body}, nil
}

func (inst *Xor) Process() {
	in1, _ := inst.ReadPinAsBool(node.In1)
	in2, _ := inst.ReadPinAsBool(node.In2)

	if (in1 && !in2) || (!in1 && in2) {
		inst.WritePinTrue(node.Out)
	} else {
		inst.WritePinFalse(node.Out)
	}
}
