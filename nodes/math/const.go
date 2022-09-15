package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Const struct {
	*node.Spec
}

func NewConst(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constNum, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Const{body}, nil
}

func (inst *Const) Process() {
	in1 := inst.ReadPin(node.In)
	inst.WritePin(node.Out, in1)
}

func (inst *Const) Cleanup() {}
