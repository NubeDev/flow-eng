package hvac

import (
	"github.com/NubeDev/flow-eng/node"
)

type DeadBand struct {
	*node.Spec
	out bool
}

func NewDeadBand(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, deadBandNode, category)
	in := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, false)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, nil, body.Inputs, false)
	deadBand := node.BuildInput(node.DeadBand, node.TypeFloat, nil, body.Inputs, false)

	inputs := node.BuildInputs(in, setPoint, deadBand)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &DeadBand{body, false}, nil
}

func (inst *DeadBand) Process() {
	input, null := inst.ReadPinAsFloat(node.Inp)
	if null {
		inst.WritePinNull(node.Out)
		return
	}
	setPoint, _ := inst.ReadPinAsFloat(node.Setpoint)
	deadBand, _ := inst.ReadPinAsFloat(node.DeadBand)
	risingEdge := setPoint + deadBand/2
	fallingEdge := deadBand - deadBand/2

	if risingEdge >= fallingEdge {
		if input <= fallingEdge {
			inst.out = false
		}
		if input >= risingEdge {
			inst.out = true
		}
	} else if risingEdge < fallingEdge {
		if input >= fallingEdge {
			inst.out = false
		}
		if input <= risingEdge {
			inst.out = true
		}
	}
	inst.WritePinBool(node.Out, inst.out)
}
