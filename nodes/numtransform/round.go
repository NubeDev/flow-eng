package numtransform

import (
	"github.com/NubeDev/flow-eng/node"
	"math"
)

type Round struct {
	*node.Spec
}

func NewRound(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, round, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false)
	decimals := node.BuildInput(node.Decimals, node.TypeFloat, 2, body.Inputs, false)
	inputs := node.BuildInputs(in, decimals)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Round{body}, nil
}

func (inst *Round) Process() {
	in, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		decimals, decNull := inst.ReadPinAsFloat(node.Decimals)
		if decNull {
			decimals = 2
		} else {
			decimals = math.Floor(decimals)
		}
		rounded := math.Round(in*math.Pow(10, decimals)) / math.Pow(10, decimals)
		inst.WritePinFloat(node.Out, rounded)
	}
}
