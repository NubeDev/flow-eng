package hvac

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type DeadBand struct {
	*node.Spec
	out bool
}

func NewDeadBand(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, deadBandNode, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, 22, body.Inputs, false, false)
	deadBand := node.BuildInput(node.DeadBand, node.TypeFloat, 1, body.Inputs, false, false)

	inputs := node.BuildInputs(in, setPoint, deadBand)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &DeadBand{body, false}, nil
}

func (inst *DeadBand) Process() {
	input, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
		return
	}
	setPoint, _ := inst.ReadPinAsFloat(node.Setpoint)
	deadBand, _ := inst.ReadPinAsFloat(node.DeadBand)
	risingEdge := setPoint + deadBand
	fallingEdge := setPoint - deadBand

	if input > risingEdge {
		inst.out = true
	}
	if input <= fallingEdge {
		inst.out = false
	}
	r := fmt.Sprint(float.RoundTo(risingEdge, 1))
	f := fmt.Sprint(float.RoundTo(fallingEdge, 1))
	msg := fmt.Sprintf("rise > %s fall <= %s", r, f)
	inst.SetSubTitle(msg)
	inst.WritePinBool(node.Out, inst.out)
}
