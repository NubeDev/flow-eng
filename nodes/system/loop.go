package system

import (
	"github.com/NubeDev/flow-eng/node"
)

type Loop struct {
	*node.Spec
}

func NewLoopCount(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, flowLoopCount, Category)
	toggleOnCount := node.BuildInput(node.TriggerOnCount, node.TypeFloat, 100, body.Inputs, false, false) // will trigger every 10 loops
	inputs := node.BuildInputs(toggleOnCount)
	outNum := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outToggle := node.BuildOutput(node.Toggle, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(outNum, outToggle)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Loop{body}, nil
}

func (inst *Loop) Process() {
	counter, _ := inst.Loop()
	toggleOnCount, null := inst.ReadPinAsUint64(node.TriggerOnCount)
	if null {
		toggleOnCount = 10
	}
	if toggleOnCount < 10 {
		toggleOnCount = 10
	}
	inst.WritePinFloat(node.Out, float64(counter))
	if toggleOnCount <= 0 {
		toggleOnCount = 1
	}
	t := counter % toggleOnCount / (toggleOnCount / 2)
	if t == 1 {
		inst.WritePinTrue(node.Toggle)

	} else {
		inst.WritePinFalse(node.Toggle)
	}
}
