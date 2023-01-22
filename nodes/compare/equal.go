package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareEqual struct {
	*node.Spec
}

func NewEqual(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, Equal, category)
	a := node.BuildInput(node.InputA, node.TypeFloat, nil, body.Inputs, nil)
	b := node.BuildInput(node.InputB, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(a, b)

	equal := node.BuildOutput(node.Equal, node.TypeBool, nil, body.Outputs)
	notEqual := node.BuildOutput(node.NotEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(equal, notEqual)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &CompareEqual{body}, nil
}

func (inst *CompareEqual) Process() {
	a, aNull := inst.ReadPinAsFloat(node.InputA)
	b, bNull := inst.ReadPinAsFloat(node.InputB)

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
