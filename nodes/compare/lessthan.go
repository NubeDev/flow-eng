package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareLessThan struct {
	*node.Spec
}

func NewLessThan(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, LessThan, category)
	a := node.BuildInput(node.InputA, node.TypeFloat, nil, body.Inputs, nil)
	b := node.BuildInput(node.InputB, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(a, b)
	graterThan := node.BuildOutput(node.LessThan, node.TypeBool, nil, body.Outputs)
	equal := node.BuildOutput(node.LessThanEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(graterThan, equal)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &CompareLessThan{body}, nil
}

func (inst *CompareLessThan) Process() {
	a, aNull := inst.ReadPinAsFloat(node.InputA)
	b, bNull := inst.ReadPinAsFloat(node.InputB)

	if aNull || bNull {
		inst.WritePin(node.LessThan, false)
		inst.WritePin(node.LessThanEqual, false)
		return
	}

	if a < b {
		inst.WritePinTrue(node.LessThan)
	} else {
		inst.WritePinFalse(node.LessThan)
	}
	if a <= b {
		inst.WritePinTrue(node.LessThanEqual)
	} else {
		inst.WritePinFalse(node.LessThanEqual)
	}
}
