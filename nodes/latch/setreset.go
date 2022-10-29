package latch

import (
	"github.com/NubeDev/flow-eng/node"
)

type SetResetLatch struct {
	*node.Spec
	currentVal bool
}

func NewSetResetLatch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, setResetLatch, category)
	set := node.BuildInput(node.Set, node.TypeBool, nil, body.Inputs)     // TODO: this input shouldn't have a manual override value
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(set, reset)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &SetResetLatch{body, false}, nil
}

func (inst *SetResetLatch) Process() {
	set, null := inst.ReadPinAsBool(node.Set)
	if null {
		inst.WritePinNull(node.Out)
		return
	}
	reset, _ := inst.ReadPinAsBool(node.Reset)
	allowResetOnSetTrue := false

	if set && !reset {
		inst.currentVal = true
	} else if allowResetOnSetTrue && reset && inst.currentVal {
		inst.currentVal = false
	} else if !set && inst.currentVal && reset {
		inst.currentVal = false
	}
	inst.WritePin(node.Out, inst.currentVal)
}
