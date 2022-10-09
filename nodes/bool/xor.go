package bool

import (
	"github.com/NubeDev/flow-eng/node"
)

type Xor struct {
	*node.Spec
}

func NewXor(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, xor, category)
	in1 := node.BuildInput(node.In1, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value
	in2 := node.BuildInput(node.In2, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(in1, in2)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Xor{body}, nil
}

func (inst *Xor) Process() {
	in1, _ := inst.ReadPinAsBool(node.In1)
	in2, _ := inst.ReadPinAsBool(node.In2)

	if (in1 && in2 != false) || (in1 != false && in2 == false) {
		inst.WritePin(node.Out, true)
	} else {
		inst.WritePin(node.Out, false)
	}
}

func (inst *Xor) Cleanup() {}
