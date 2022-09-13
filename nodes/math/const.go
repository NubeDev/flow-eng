package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Const struct {
	*node.BaseNode
}

func NewConst(body *node.BaseNode) (node.Node, error) {
	body = node.Defaults(body, constNum, category)
	inputs := node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Const{body}, nil
}

func (inst *Const) Process() {
	in1 := inst.ReadPin(node.In1)
	inst.WritePin(node.Out1, in1)
}

func (inst *Const) Cleanup() {}
