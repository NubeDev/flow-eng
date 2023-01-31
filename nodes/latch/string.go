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
	input := node.BuildInput(node.Inp, node.TypeString, nil, body.Inputs, nil)
	latch := node.BuildInput(node.Latch, node.TypeBool, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(input, latch)
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeString, "", body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &StringLatch{body, "", false}, nil
}

func (inst *StringLatch) Process() {
	input, _ := inst.ReadPinAsString(node.Inp)
	latch, _ := inst.ReadPinAsBool(node.Latch)
	latchBool := latch

	if latchBool && !inst.lastTrigger {
		inst.currentVal = input
	}
	inst.lastTrigger = latchBool

	inst.WritePin(node.Outp, inst.currentVal)
}
