package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type Between struct {
	*node.Spec
}

func NewBetween(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logicCompare, category)
	var inNames = []string{InputName(node.In), InputName(node.From), InputName(node.To)}
	var outNames = []string{OutputName(node.Out), OutputName(node.OutNot), OutputName(node.Above), OutputName(node.Below)}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 3, 3, 3, body.Inputs, inNames)...)
	outputs := node.BuildOutputs(node.DynamicOutputs(node.TypeFloat, nil, 4, 4, 4, body.Outputs, outNames)...)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Between{body}, nil
}

func (inst *Between) Process() {
	Process(inst)
}

func (inst *Between) Cleanup() {}
