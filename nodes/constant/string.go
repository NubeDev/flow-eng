package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type ConstString struct {
	*node.Spec
}

func NewString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constStr, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(constHelp)
	body.SetAllowPayload()
	body.SetPayloadType(node.TypeString)
	return &ConstString{body}, nil
}

func (inst *ConstString) Process() {
	// context menu payload overrides
	v, null := inst.ReadPayloadAsString2()
	if !null {
		inst.OverrideInputValue(node.In, v)
	}

	in1 := inst.ReadPin(node.In)
	inst.WritePin(node.Out, in1)
}
