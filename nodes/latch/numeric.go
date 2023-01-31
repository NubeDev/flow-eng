package latch

import (
	"github.com/NubeDev/flow-eng/node"
)

type NumLatch struct {
	*node.Spec
	currentVal  float64
	lastTrigger bool
}

func NewNumLatch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, numLatch, category)
	input := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)
	latch := node.BuildInput(node.Latch, node.TypeBool, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(input, latch)

	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &NumLatch{body, 0, false}, nil
}

func (inst *NumLatch) Process() {
	input, _ := inst.ReadPinAsFloat(node.Inp)
	latch, _ := inst.ReadPinAsBool(node.Latch)
	latchBool := latch

	if latchBool && !inst.lastTrigger {
		inst.currentVal = input
	}
	inst.lastTrigger = latchBool

	inst.WritePinFloat(node.Outp, inst.currentVal)
}
