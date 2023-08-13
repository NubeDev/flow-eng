package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constNum, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(constHelp)
	body.SetAllowPayload()
	body.SetPayloadType(node.TypeFloat)
	return &Number{body}, nil
}

func (inst *Number) Process() {
	// context menu payload overrides
	v, noPayload, nullPayload := inst.ReadPayloadAsFloat()
	if !noPayload && !nullPayload {
		inst.OverrideInputValue(node.In, v)
	} else if nullPayload {
		inst.WritePinNull(node.Out)
		return
	}
	in1, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinFloat(node.Out, in1)
	}

}
