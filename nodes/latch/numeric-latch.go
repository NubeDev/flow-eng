package latch

import (
	"github.com/NubeDev/flow-eng/node"
)

type NumLatch struct {
	*node.Spec
	lastTrigger bool
}

const (
	category = "latch"
	numLatch = "numeric-latch"
)

func NewNumLatch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, numLatch, category)
	input := node.BuildInput(node.Input_, node.TypeFloat, nil, body.Inputs)
	latch := node.BuildInput(node.Latch, node.TypeFloat, nil, body.Inputs)

	inputs := node.BuildInputs(input, latch)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &NumLatch{body, false}, nil
}

func (inst *NumLatch) Process() {
	input := inst.ReadPinAsFloat(node.Input_)
	latch := inst.ReadPinAsFloat(node.Latch)
	latchBool := latch == 1

	if latchBool && !inst.lastTrigger {
		inst.WritePin(node.Out, input)
	}
	inst.lastTrigger = latchBool
}

func (inst *NumLatch) Cleanup() {}
