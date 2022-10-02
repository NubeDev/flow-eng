package bool

import (
	"github.com/NubeDev/flow-eng/node"
)

type Not struct {
	*node.Spec
}

func NewNot(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, not, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(in)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Not{body}, nil
}

func (inst *Not) Process() {
	in := inst.ReadPinBool(node.In)
	if in {
		inst.WritePin(node.Out, false)
	} else {
		inst.WritePin(node.Out, true)
	}
}

func (inst *Not) Cleanup() {}
