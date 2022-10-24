package system

import (
	"github.com/NubeDev/flow-eng/node"
)

type Loop struct {
	*node.Spec
	toggle bool
}

func NewLoopCount(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowLoopCount, category)
	toggleOnCount := node.BuildInput(node.TriggerOnCount, node.TypeFloat, 10, body.Inputs) // will trigger every 10 loops
	inputs := node.BuildInputs(toggleOnCount)
	outNum := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outToggle := node.BuildOutput(node.Toggle, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(outNum, outToggle)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Loop{body, false}, nil
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
			inst.WritePinTrue(node.Toggle)
			inst.toggle = false
		} else {
			inst.WritePinFalse(node.Toggle)
			inst.toggle = true
		}
		return
	}
	if counter%toggleOnCount == 0 {
		inst.WritePinTrue(node.Toggle)
	} else {
		inst.WritePinFalse(node.Toggle)
	}
}
