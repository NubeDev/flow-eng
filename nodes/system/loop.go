package system

import (
	"github.com/NubeDev/flow-eng/node"
)

type Loop struct {
	*node.Spec
}

func NewLoopCount(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowLoopCount, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(in)
	outNum := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outToggle := node.BuildOutput(node.Toggle, node.TypeFloat, nil, body.Outputs)

	outputs := node.BuildOutputs(outNum, outToggle)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Loop{body}, nil
}

var counter uint64

func (inst *Loop) Process() {
	counter++
	inst.WritePin(node.Out, counter)
	if counter%2 == 0 {
		inst.WritePin(node.Toggle, 0)
	} else {
		inst.WritePin(node.Toggle, 1)
	}
}

func (inst *Loop) Cleanup() {}
