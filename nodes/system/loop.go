package system

import (
	"github.com/NubeDev/flow-eng/node"
)

type Loop struct {
	*node.Spec
}

func NewLoopCount(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowLoopCount, category)
	toggleOnCount := node.BuildInput(node.Count, node.TypeFloat, 10, body.Inputs) // will trigger every 10 loops
	inputs := node.BuildInputs(toggleOnCount)
	outNum := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outToggle := node.BuildOutput(node.Toggle, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(outNum, outToggle)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Loop{body}, nil
}

var counter uint64

func (inst *Loop) Process() {
	counter++
	toggleOnCount, _ := inst.ReadPinAsUint64(node.Count)
	if toggleOnCount <= 0 {
		toggleOnCount = 2
	}
	inst.WritePin(node.Out, counter)
	if counter%toggleOnCount == 0 {
		inst.WritePinTrue(node.Toggle)
	} else {
		inst.WritePinFalse(node.Toggle)
	}
}

func (inst *Loop) Cleanup() {}
