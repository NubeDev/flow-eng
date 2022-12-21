package bool

import (
	"github.com/NubeDev/flow-eng/node"
)

type Toggle struct {
	*node.Spec
	currentOut bool
	lastIn     bool
}

func NewToggle(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, toggle, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(in)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Toggle{body, false, false}, nil
}

func (inst *Toggle) Process() {
	in, _ := inst.ReadPinAsBool(node.In)
	inAsBool := in == true
	if inAsBool && !inst.lastIn {
		inst.currentOut = !inst.currentOut
	}
	inst.lastIn = inAsBool

	if inst.currentOut {
		inst.WritePin(node.Out, true)
	} else {
		inst.WritePin(node.Out, false)
	}
}
