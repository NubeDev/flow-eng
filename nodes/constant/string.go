package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type ConstString struct {
	*node.Spec
}

func NewString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constStr, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(constHelp)
	return &ConstString{body}, nil
}

func (inst *ConstString) Process() {
	in1 := inst.ReadPin(node.In)
	inst.WritePin(node.Out, in1)
}
