package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareGreaterThan struct {
	*node.Spec
}

func NewGreaterThan(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, GreaterThan, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs, nil)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in1, in2)
	graterThan := node.BuildOutput(node.GreaterThan, node.TypeBool, nil, body.Outputs)
	equal := node.BuildOutput(node.GreaterThanEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(graterThan, equal)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &CompareGreaterThan{body}, nil
}

func (inst *CompareGreaterThan) Process() {
	a, aNull := inst.ReadPinAsFloat(node.In1)
	b, bNull := inst.ReadPinAsFloat(node.In2)

	if aNull || bNull {
		inst.WritePinBool(node.GreaterThan, false)
		inst.WritePinBool(node.GreaterThanEqual, false)
		return
	}

	if a > b {
		inst.WritePinTrue(node.GreaterThan)
	} else {
		inst.WritePinFalse(node.GreaterThan)
	}
	if a >= b {
		inst.WritePinTrue(node.GreaterThanEqual)
	} else {
		inst.WritePinFalse(node.GreaterThanEqual)
	}
}
