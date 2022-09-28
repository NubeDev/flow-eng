package bool

import (
	"github.com/NubeDev/flow-eng/node"
)

type Not struct {
	*node.Spec
}

func NewNot(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, not, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(in)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Not{body}, nil
}

func (inst *Not) Process() {
	in := inst.ReadPinAsFloat(node.In)
	if in == 1 {
		inst.WritePin(node.Out, 0)
	} else {
		inst.WritePin(node.Out, 1)
	}
}

func (inst *Not) Cleanup() {}
