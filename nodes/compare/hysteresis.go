package compare

import (
	"github.com/NubeDev/flow-eng/node"
)

type Hysteresis struct {
	*node.Spec
}

func NewHysteresis(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, hysteresis, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	risingEdge := node.BuildInput(node.RisingEdge, node.TypeFloat, 20, body.Inputs)
	fallingEdge := node.BuildInput(node.FallingEdge, node.TypeFloat, 10, body.Inputs)

	inputs := node.BuildInputs(in, risingEdge, fallingEdge)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Hysteresis{body}, nil
}

func (inst *Hysteresis) Process() {
	value := inst.ReadPinAsFloat(node.In)
	var out bool
	risingEdge := inst.ReadPinAsFloat(node.RisingEdge)
	fallingEdge := inst.ReadPinAsFloat(node.FallingEdge)

	if risingEdge > fallingEdge {
		if value <= fallingEdge {
			out = false
		}
		if value >= risingEdge {
			out = true
		}
	} else if risingEdge < fallingEdge {
		if value >= fallingEdge {
			out = false
		}
		if value <= risingEdge {
			out = true
		}
	}
	if out {
		inst.WritePin(node.Out, 1)
	} else {
		inst.WritePin(node.Out, 0)
	}

}

func (inst *Hysteresis) Cleanup() {}