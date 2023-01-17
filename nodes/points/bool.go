package point

import (
	"github.com/NubeDev/flow-eng/node"
)

type Boolean struct {
	*node.Spec
}

func NewBoolean(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pointBoolean, category)
	in1 := node.BuildInput(node.In1, node.TypeBool, nil, body.Inputs, nil)
	in2 := node.BuildInput(node.In2, node.TypeBool, nil, body.Inputs, nil)
	in3 := node.BuildInput(node.In3, node.TypeBool, nil, body.Inputs, nil)
	in4 := node.BuildInput(node.In4, node.TypeBool, nil, body.Inputs, nil)
	body.Inputs = node.BuildInputs(in1, in2, in3, in4)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &Boolean{body}, nil
}

func (inst *Boolean) Process() {
	in1, in1Null := inst.ReadPinAsBool(node.In1)
	in2, in2Null := inst.ReadPinAsBool(node.In2)
	in3, in3Null := inst.ReadPinAsBool(node.In3)
	in4, in4Null := inst.ReadPinAsBool(node.In4)
	if !in1Null {
		inst.WritePin(node.Out, in1)
		return
	}
	if !in2Null {
		inst.WritePin(node.Out, in2)
		return
	}
	if !in3Null {
		inst.WritePin(node.Out, in3)
		return
	}
	if !in4Null {
		inst.WritePin(node.Out, in4)
		return
	}
	inst.WritePinNull(node.Out)
}
