package latch

import (
	"github.com/NubeDev/flow-eng/node"
)

type StringLatch struct {
	*node.Spec
	currentVal  string
	lastTrigger bool
}

func NewStringLatch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, stringLatch, category)
	input := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	latch := node.BuildInput(node.Latch, node.TypeFloat, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(input, latch)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, "", body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &StringLatch{body, "", false}, nil
}

func (inst *StringLatch) Process() {
	input := inst.ReadPinAsString(node.In)
	latch, _ := inst.ReadPinAsFloat(node.Latch)
	latchBool := latch == 1

	if latchBool && !inst.lastTrigger {
		inst.currentVal = input
	}
	inst.lastTrigger = latchBool

	inst.WritePin(node.Out, inst.currentVal)
}

func (inst *StringLatch) Cleanup() {}
