package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type Boolean struct {
	*node.Spec
}

func NewBoolean(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constBool, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(constHelp)
	body.SetAllowPayload()
	body.SetPayloadType(node.TypeBool)
	return &Boolean{body}, nil
}

func (inst *Boolean) Process() {
	// context menu payload overrides
	v, null := inst.ReadPayloadAsBool()
	if !null {
		inst.OverrideInputValue(node.In, v)
	}

	in, null := inst.ReadPinAsBool(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinBool(node.Out, in)
	}
}
