package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type Boolean struct {
	*node.Spec
}

func NewBoolean(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constBool, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeBool, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(constHelp)
	return &Boolean{body}, nil
}

func (inst *Boolean) Process() {
	v, null := inst.ReadPinAsBool(node.Inp)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinBool(node.Out, v)
	}
}
