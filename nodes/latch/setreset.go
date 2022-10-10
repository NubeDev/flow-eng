package latch

import (
	"github.com/NubeDev/flow-eng/node"
)

type SetResetLatch struct {
	*node.Spec
	currentVal float64
}

func NewSetResetLatch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, setResetLatch, category)
	set := node.BuildInput(node.Set, node.TypeFloat, nil, body.Inputs)     // TODO: this input shouldn't have a manual override value
	reset := node.BuildInput(node.Reset, node.TypeFloat, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(set, reset)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &SetResetLatch{body, 0}, nil
}

func (inst *SetResetLatch) Process() {
	set, _ := inst.ReadPinAsFloat(node.Set)
	reset, _ := inst.ReadPinAsFloat(node.Reset)
	allowResetOnSetTrue := false

	if set == 1 && reset != 1 {
		inst.currentVal = 1
	} else if allowResetOnSetTrue && reset == 1 && inst.currentVal == 1 {
		inst.currentVal = 0
	} else if set != 1 && inst.currentVal == 1 && reset == 1 {
		inst.currentVal = 0
	}

	inst.WritePin(node.Out, inst.currentVal)
}
