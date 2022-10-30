package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type CompareGreater struct {
	*node.Spec
}

func NewCompareGreater(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logicCompareGreater, category)
	a := node.BuildInput(node.InputA, node.TypeFloat, nil, body.Inputs)
	b := node.BuildInput(node.InputB, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(a, b)
	graterThan := node.BuildOutput(node.GreaterThan, node.TypeBool, nil, body.Outputs)
	equal := node.BuildOutput(node.GreaterThanEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(graterThan, equal)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &CompareGreater{body}, nil
}

func (inst *CompareGreater) Process() {
	a, _ := inst.ReadPinAsFloat(node.InputA)
	b, _ := inst.ReadPinAsFloat(node.InputB)

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

type CompareLess struct {
	*node.Spec
}

func NewCompareLess(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logicCompareLess, category)
	a := node.BuildInput(node.InputA, node.TypeFloat, nil, body.Inputs)
	b := node.BuildInput(node.InputB, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(a, b)
	lessThan := node.BuildOutput(node.LessThan, node.TypeBool, nil, body.Outputs)
	equal := node.BuildOutput(node.LessThanEqual, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(lessThan, equal)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &CompareLess{body}, nil
}

func (inst *CompareLess) Process() {
	a, _ := inst.ReadPinAsFloat(node.InputA)
	b, _ := inst.ReadPinAsFloat(node.InputB)

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

type CompareEqual struct {
	*node.Spec
}

func NewCompareEqual(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logicCompareEqual, category)
	a := node.BuildInput(node.InputA, node.TypeFloat, nil, body.Inputs)
	b := node.BuildInput(node.InputB, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(a, b)
	equal := node.BuildOutput(node.Equal, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(equal)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &CompareEqual{body}, nil
}

func (inst *CompareEqual) Process() {
	a, _ := inst.ReadPinAsFloat(node.InputA)
	b, _ := inst.ReadPinAsFloat(node.InputB)
	if a == b {
		inst.WritePinTrue(node.Equal)
	} else {
		inst.WritePinFalse(node.Equal)
	}
}
