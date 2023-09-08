package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type Hysteresis struct {
	*node.Spec
	currentVal bool
}

func NewHysteresis(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, hysteresis, Category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	risingEdge := node.BuildInput(node.RisingEdge, node.TypeFloat, 100, body.Inputs, false, false)
	fallingEdge := node.BuildInput(node.FallingEdge, node.TypeFloat, 0, body.Inputs, false, false)
	inputs := node.BuildInputs(in, risingEdge, fallingEdge)

	output := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outNot := node.BuildOutput(node.OutNot, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(output, outNot)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Hysteresis{body, false}, nil
}

func (inst *Hysteresis) Process() {
	value, inNull := inst.ReadPinAsFloat(node.In)
	risingEdge, riseNull := inst.ReadPinAsFloat(node.RisingEdge)
	fallingEdge, fallNull := inst.ReadPinAsFloat(node.FallingEdge)

	if riseNull || fallNull || inNull {
		inst.WritePinFalse(node.Out)
		inst.WritePinTrue(node.Out)
		inst.currentVal = false
		return
	}

	if risingEdge > fallingEdge {
		if value <= fallingEdge {
			inst.currentVal = false
		}
		if value >= risingEdge {
			inst.currentVal = true
		}
	} else if risingEdge < fallingEdge {
		if value >= fallingEdge {
			inst.currentVal = false
		}
		if value <= risingEdge {
			inst.currentVal = false
		}
	} else if risingEdge == fallingEdge {
		inst.currentVal = value > risingEdge
	}

	inst.WritePinBool(node.Out, inst.currentVal)
	inst.WritePinBool(node.OutNot, !inst.currentVal)
}
