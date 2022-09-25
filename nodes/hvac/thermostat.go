package hvac

import "github.com/NubeDev/flow-eng/node"

type Thermostat struct {
	*node.Spec
}

// input
// enable
// sp
// cool offset
// heat offset

func NewThermostat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, deadBandNode, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	setPoint := node.BuildInput(node.SetPoint, node.TypeFloat, nil, body.Inputs)
	deadBand := node.BuildInput(node.DeadBand, node.TypeFloat, nil, body.Inputs)

	inputs := node.BuildInputs(in, setPoint, deadBand)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Thermostat{body}, nil
}

func (inst *Thermostat) Process() {
	input := inst.ReadPinAsFloat(node.In)
	var out bool
	setPoint := inst.ReadPinAsFloat(node.SetPoint)
	deadBand := inst.ReadPinAsFloat(node.DeadBand)
	risingEdge := setPoint + deadBand/2
	fallingEdge := deadBand - deadBand/2

	if risingEdge >= fallingEdge {
		if input <= fallingEdge {
			out = false
		}
		if input >= risingEdge {
			out = true
		}
	} else if risingEdge < fallingEdge {
		if input >= fallingEdge {
			out = false
		}
		if input <= risingEdge {
			out = true
		}
	}
	if out {
		inst.WritePin(node.Out, 1)
	} else {
		inst.WritePin(node.Out, 0)
	}

}

func (inst *Thermostat) Cleanup() {}
