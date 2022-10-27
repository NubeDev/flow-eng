package system

import (
	"github.com/NubeDev/flow-eng/node"
)

type Loop struct {
	*node.Spec
	toggle bool
	count  int
}

func NewLoopCount(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowLoopCount, category)
	toggleOnCount := node.BuildInput(node.TriggerOnCount, node.TypeFloat, 10, body.Inputs) // will trigger every 10 loops
	inputs := node.BuildInputs(toggleOnCount)
	outNum := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outToggle := node.BuildOutput(node.Trigger, node.TypeBool, nil, body.Outputs)
	outToggle2 := node.BuildOutput(node.Toggle, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(outNum, outToggle, outToggle2)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Loop{body, false, 0}, nil
}

func (inst *Loop) Process() {
	counter, _ := inst.Loop()
	toggleOnCount, _ := inst.ReadPinAsUint64(node.TriggerOnCount)
	inst.WritePin(node.Out, counter)
	if toggleOnCount <= 0 {
		toggleOnCount = 1
	}
	if toggleOnCount == 1 {
		if inst.toggle {
			inst.WritePinTrue(node.Trigger)
			inst.toggle = false
		} else {
			inst.WritePinFalse(node.Trigger)
			inst.toggle = true
		}
		return
	}
	t := counter % toggleOnCount / (toggleOnCount / 2)
	if t == 1 {
		inst.WritePinTrue(node.Toggle)
	} else {
		inst.WritePinFalse(node.Toggle)
	}
	if counter%toggleOnCount == 0 {
		inst.WritePinTrue(node.Trigger)
	} else {
		inst.WritePinFalse(node.Trigger)
	}
}
