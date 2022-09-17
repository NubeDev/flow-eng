package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type Compare struct {
	*node.Spec
}

func NewCompare(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logicCompare, category)
	var names = []string{node.SetOutputName(node.GraterThan), node.SetOutputName(node.LessThan), node.SetOutputName(node.Equal)}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 2, 2, 2, body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.DynamicOutputs(node.TypeFloat, nil, 3, 3, 3, body.Outputs, names)...)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Compare{body}, nil
}

func (inst *Compare) Process() {
	Process(inst)
}

func (inst *Compare) Cleanup() {}
