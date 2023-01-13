package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareEqualString struct {
	*node.Spec
}

func NewEqualString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, EqualString, category)
	a := node.BuildInput(node.InputA, node.TypeString, nil, body.Inputs, nil)
	b := node.BuildInput(node.InputB, node.TypeString, nil, body.Inputs, nil)
	inputs := node.BuildInputs(a, b)

	equal := node.BuildOutput(node.Equal, node.TypeBool, nil, body.Outputs)
	notEqual := node.BuildOutput(node.NotEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(equal, notEqual)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &CompareEqualString{body}, nil
}

func (inst *CompareEqualString) Process() {
	a, aNull := inst.ReadPinAsString(node.InputA)
	b, bNull := inst.ReadPinAsString(node.InputB)

	if (aNull && bNull) || a == b {
		inst.WritePin(node.Equal, true)
		inst.WritePin(node.NotEqual, false)
		return
	} else {
		inst.WritePin(node.Equal, false)
		inst.WritePin(node.NotEqual, true)
		return
	}

}
