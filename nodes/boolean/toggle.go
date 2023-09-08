package boolean

import (
	"github.com/NubeDev/flow-eng/node"
)

type Toggle struct {
	*node.Spec
	currentOut bool
	lastIn     bool
}

func NewToggle(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, toggle, Category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false, true)
	inputs := node.BuildInputs(in)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Toggle{body, false, true}, nil
}

func (inst *Toggle) Process() {
	resetOnNullOrDisconnect := false
	in, null := inst.ReadPinAsBool(node.In)
	if null && resetOnNullOrDisconnect {
		inst.WritePinFalse(node.Out)
		inst.currentOut = false
		return
	}
	if !null {
		if in && !inst.lastIn {
			inst.currentOut = !inst.currentOut
		}
		inst.lastIn = in
	}
	inst.WritePinBool(node.Out, inst.currentOut)
}
