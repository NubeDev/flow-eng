package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type Boolean struct {
	*node.Spec
}

func NewBoolean(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, constBool, Category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(constHelp)
	body.SetAllowPayload()
	body.SetPayloadType(node.TypeBool)
	return &Boolean{body}, nil
}

func (inst *Boolean) Process() {
	// context menu payload overrides
	v, noPayload, nullPayload := inst.ReadPayloadAsBool()
	if !noPayload && !nullPayload {
		inst.OverrideInputValue(node.In, v)
	} else if nullPayload {
		inst.WritePinNull(node.Out)
		return
	}

	if !noPayload && !nullPayload {
		inst.OverrideInputValue(node.In, v)
	} else if nullPayload {
		inst.WritePinNull(node.Out)
		return
	}

	in, null := inst.ReadPinAsBool(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinBool(node.Out, in)
	}
}
