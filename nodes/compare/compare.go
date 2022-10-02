package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type Compare struct {
	*node.Spec
}

func NewCompare(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logicCompare, category)
	a := node.BuildInput(node.InputA, node.TypeFloat, nil, body.Inputs)
	b := node.BuildInput(node.InputB, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(a, b)
	graterThan := node.BuildOutput(node.GraterThan, node.TypeBool, nil, body.Outputs)
	lessThan := node.BuildOutput(node.LessThan, node.TypeBool, nil, body.Outputs)
	equal := node.BuildOutput(node.Equal, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(graterThan, lessThan, equal)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Compare{body}, nil
}

func (inst *Compare) Process() {
	a := inst.ReadPinAsFloat(node.InputA)
	b := inst.ReadPinAsFloat(node.InputB)

	if a > b {
		inst.WritePin(node.GraterThan, true)
	} else {
		inst.WritePin(node.GraterThan, false)
	}
	if a < b {
		inst.WritePin(node.LessThan, true)
	} else {
		inst.WritePin(node.LessThan, false)
	}
	if a == b {
		inst.WritePin(node.Equal, true)
	} else {
		inst.WritePin(node.Equal, false)
	}
}

func (inst *Compare) Cleanup() {}
