package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareLessThan struct {
	*node.Spec
}

func NewLessThan(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, LessThan, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs, false, false)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in1, in2)
	greaterThan := node.BuildOutput(node.LessThan, node.TypeBool, nil, body.Outputs)
	equal := node.BuildOutput(node.LessThanEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(greaterThan, equal)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &CompareLessThan{body}, nil
}

func (inst *CompareLessThan) Process() {
	a, aNull := inst.ReadPinAsFloat(node.In1)
	b, bNull := inst.ReadPinAsFloat(node.In2)

	if aNull || bNull {
		inst.WritePinBool(node.LessThan, false)
		inst.WritePinBool(node.LessThanEqual, false)
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
