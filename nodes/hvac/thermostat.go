package hvac

import "github.com/NubeDev/flow-eng/node"

type Thermostat struct {
	*node.Spec
	out bool
}

// input
// enable
// sp
// cool offset
// heat offset

func NewThermostat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, deadBandNode, Category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, 22, body.Inputs, false, false)
	deadBand := node.BuildInput(node.DeadBand, node.TypeFloat, 2, body.Inputs, false, false)

	inputs := node.BuildInputs(in, setPoint, deadBand)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Thermostat{body, false}, nil
}

func (inst *Thermostat) Process() {
	input, null := inst.ReadPinAsFloat(node.In)
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
