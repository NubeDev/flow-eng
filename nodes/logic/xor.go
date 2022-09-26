package logic

import (
	"github.com/NubeDev/flow-eng/node"
)

type Xor struct {
	*node.Spec
}

func NewXor(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, xor, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs)

	inputs := node.BuildInputs(in1, in2)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Xor{body}, nil
}

func (inst *Xor) Process() {
	in1 := inst.ReadPinAsFloat(node.In1)
	in2 := inst.ReadPinAsFloat(node.In2)

	if (in1 == 0 && in2 != 0) || (in1 != 0 && in2 == 0) {
		inst.WritePin(node.Out, 1)
	} else {
		inst.WritePin(node.Out, 0)
	}
}

func (inst *Xor) Cleanup() {}
