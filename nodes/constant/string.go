package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type ConstString struct {
	*node.Spec
}

func NewString(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, constStr, Category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(constHelp)
	body.SetAllowPayload()
	body.SetPayloadType(node.TypeString)
	return &ConstString{body}, nil
}

func (inst *ConstString) Process() {
	// context menu payload overrides
	v, noPayload, nullPayload := inst.ReadPayloadAsString()
	if !noPayload && !nullPayload {
		inst.OverrideInputValue(node.In, v)
	} else if nullPayload {
		inst.WritePinNull(node.Out)
		return
	}

	in1 := inst.ReadPin(node.In)
	if in1 == nil || in1 == "null" || in1 == "" {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePin(node.Out, in1)
	}
}
