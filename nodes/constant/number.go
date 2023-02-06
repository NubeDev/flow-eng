package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constNum, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(constHelp)
	body.SetAllowPayload()
	body.SetPayloadType(node.TypeNumber)
	return &Number{body}, nil
}

func (inst *Number) Process() {
	in1, null := inst.ReadPinAsFloat(node.Inp)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinFloat(node.Out, in1)
	}
	/*
		v, null := inst.ReadPayloadAsFloat()
		if !null {
			inst.OverrideInputValue(node.Inp, v)
		}
	*/
}
