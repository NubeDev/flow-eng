package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareEqualString struct {
	*node.Spec
}

func NewEqualString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, EqualString, category)
	in1 := node.BuildInput(node.In1, node.TypeString, nil, body.Inputs, nil)
	in2 := node.BuildInput(node.In2, node.TypeString, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in1, in2)

	equal := node.BuildOutput(node.Equal, node.TypeBool, nil, body.Outputs)
	notEqual := node.BuildOutput(node.NotEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(equal, notEqual)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &CompareEqualString{body}, nil
}

func (inst *CompareEqualString) Process() {
	a, aNull := inst.ReadPinAsString(node.In1)
	b, bNull := inst.ReadPinAsString(node.In2)

	if (aNull && bNull) || a == b {
		inst.WritePinBool(node.Equal, true)
		inst.WritePinBool(node.NotEqual, false)
		return
	} else {
		inst.WritePinBool(node.Equal, false)
		inst.WritePinBool(node.NotEqual, true)
		return
	}

}
