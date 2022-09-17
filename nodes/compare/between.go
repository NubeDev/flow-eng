package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type Between struct {
	*node.Spec
}

func NewBetween(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, between, category)
	var inNames = []string{node.SetInputName(node.In), node.SetInputName(node.From), node.SetInputName(node.To)}
	var outNames = []string{node.SetOutputName(node.Out), node.SetOutputName(node.OutNot), node.SetOutputName(node.Above), node.SetOutputName(node.Below)}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 3, 0, 0, body.Inputs, inNames)...)
	outputs := node.BuildOutputs(node.DynamicOutputs(node.TypeFloat, nil, 4, 0, 0, body.Outputs, outNames)...)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Between{body}, nil
}

func (inst *Between) Process() {
	Process(inst)
}

func (inst *Between) Cleanup() {}
