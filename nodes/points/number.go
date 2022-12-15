package point

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/enescakir/emoji"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pointNumber, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs)
	in3 := node.BuildInput(node.In3, node.TypeFloat, nil, body.Inputs)
	in4 := node.BuildInput(node.In4, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in1, in2, in3, in4)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	body.SetAllowPayload()
	body.SetPayloadType(node.TypeNumber)
	body.SetIcon(string(emoji.Label))
	return &Number{body}, nil
}

func (inst *Number) Process() {
	in1, in1Null := inst.ReadPinAsFloat(node.In1)
	in2, in2Null := inst.ReadPinAsFloat(node.In2)
	in3, in3Null := inst.ReadPinAsFloat(node.In3)
	in4, in4Null := inst.ReadPinAsFloat(node.In4)
	if !in1Null {
		inst.WritePinFloat(node.Out, in1)
		return
	}
	if !in2Null {
		inst.WritePinFloat(node.Out, in2)
		return
	}
	if !in3Null {
		inst.WritePinFloat(node.Out, in3)
		return
	}
	if !in4Null {
		inst.WritePinFloat(node.Out, in4)
		return
	}
	inst.WritePinNull(node.Out)
}
