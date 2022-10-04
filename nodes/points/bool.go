package point

import (
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/NubeDev/flow-eng/node"
)

type Boolean struct {
	*node.Spec
}

func NewBoolean(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pointBoolean, category)
	in1 := node.BuildInput(node.In1, node.TypeBool, nil, body.Inputs)
	in2 := node.BuildInput(node.In2, node.TypeBool, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in1, in2)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &Boolean{body}, nil
}

func (inst *Boolean) Process() {
	in1 := inst.ReadPinBoolPointer(node.In1)
	in2 := inst.ReadPinBoolPointer(node.In2)
	if in1 != nil {
		inst.WritePin(node.Out, boolean.NonNil(in1))
	} else {
		inst.WritePin(node.Out, boolean.NonNil(in2))
	}
}
func (inst *Boolean) Cleanup() {}
