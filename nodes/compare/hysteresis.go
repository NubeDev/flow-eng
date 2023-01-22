package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type Hysteresis struct {
	*node.Spec
	currentVal bool
}

func NewHysteresis(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, hysteresis, category)
	in := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)
	risingEdge := node.BuildInput(node.RisingEdge, node.TypeFloat, 20, body.Inputs, nil)
	fallingEdge := node.BuildInput(node.FallingEdge, node.TypeFloat, 10, body.Inputs, nil)
	inputs := node.BuildInputs(in, risingEdge, fallingEdge)

	output := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	outNot := node.BuildOutput(node.OutNot, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(output, outNot)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Hysteresis{body, false}, nil
}

func (inst *Hysteresis) Process() {
	value, inNull := inst.ReadPinAsFloat(node.Inp)
	risingEdge, riseNull := inst.ReadPinAsFloat(node.RisingEdge)
	fallingEdge, fallNull := inst.ReadPinAsFloat(node.FallingEdge)

	if riseNull || fallNull || inNull {
		inst.WritePinFalse(node.Outp)
		inst.WritePinTrue(node.Outp)
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

	inst.WritePinBool(node.Outp, inst.currentVal)
	inst.WritePinBool(node.OutNot, !inst.currentVal)
}
