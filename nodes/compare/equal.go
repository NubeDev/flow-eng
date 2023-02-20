package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareEqual struct {
	*node.Spec
}

func NewEqual(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, Equal, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs, false, false)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in1, in2)

	equal := node.BuildOutput(node.Equal, node.TypeBool, nil, body.Outputs)
	notEqual := node.BuildOutput(node.NotEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(equal, notEqual)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &CompareEqual{body}, nil
}

func (inst *CompareEqual) Process() {
	a, aNull := inst.ReadPinAsFloat(node.In1)
	b, bNull := inst.ReadPinAsFloat(node.In2)

	if (aNull && bNull) || (!aNull && !bNull && a == b) {
		inst.WritePinBool(node.Equal, true)
		inst.WritePinBool(node.NotEqual, false)
		return
	} else {
		inst.WritePinBool(node.Equal, false)
		inst.WritePinBool(node.NotEqual, true)
		return
	}

}
