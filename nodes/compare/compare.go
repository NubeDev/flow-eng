package compare

import (
	"github.com/NubeDev/flow-eng/node"
	"strings"
)

type Compare struct {
	*node.Spec
}

func InputName(name node.InputName) string {
	return strings.ToLower(string(name))
}

func OutputName(name node.OutputName) string {
	return strings.ToLower(string(name))
}

func NewCompare(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logicCompare, category)
	var names = []string{OutputName(node.GraterThan), OutputName(node.LessThan), OutputName(node.Equal)}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 2, 2, 2, body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.DynamicOutputs(node.TypeFloat, nil, 3, 3, 3, body.Outputs, names)...)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Compare{body}, nil
}

func (inst *Compare) Process() {
	Process(inst)
}

func (inst *Compare) Cleanup() {}
