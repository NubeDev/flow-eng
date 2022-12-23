package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareGreaterThan struct {
	*node.Spec
}

func NewGreaterThan(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, GreaterThan, category)
	a := node.BuildInput(node.InputA, node.TypeFloat, nil, body.Inputs)
	b := node.BuildInput(node.InputB, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(a, b)
	graterThan := node.BuildOutput(node.GreaterThan, node.TypeBool, nil, body.Outputs)
	equal := node.BuildOutput(node.GreaterThanEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(graterThan, equal)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &CompareGreaterThan{body}, nil
}

func (inst *CompareGreaterThan) Process() {
	a, aNull := inst.ReadPinAsFloat(node.InputA)
	b, bNull := inst.ReadPinAsFloat(node.InputB)

	if aNull || bNull {
		inst.WritePin(node.GreaterThan, false)
		inst.WritePin(node.GreaterThanEqual, false)
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
