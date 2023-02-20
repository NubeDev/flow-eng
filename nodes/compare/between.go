package compare

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
)

type Between struct {
	*node.Spec
}

func NewBetween(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, betweenNode, category)

	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	from := node.BuildInput(node.From, node.TypeFloat, 100, body.Inputs, false, false)
	to := node.BuildInput(node.To, node.TypeFloat, 0, body.Inputs, false, false)
	inputs := node.BuildInputs(in, from, to)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outNot := node.BuildOutput(node.OutNot, node.TypeBool, nil, body.Outputs)
	above := node.BuildOutput(node.Above, node.TypeBool, nil, body.Outputs)
	below := node.BuildOutput(node.Below, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out, outNot, above, below)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Between{body}, nil
}

func (inst *Between) Process() {
	in, inNull := inst.ReadPinAsFloat(node.In)
	from, fromNull := inst.ReadPinAsFloat(node.From)
	to, toNull := inst.ReadPinAsFloat(node.To)

	if inNull || fromNull || toNull {
		inst.WritePinBool(node.Out, false)
		inst.WritePinBool(node.OutNot, true)
		inst.WritePinBool(node.Above, false)
		inst.WritePinBool(node.Below, false)
	}

	between, below, above := array.Between(in, from, to)
	inst.WritePinBool(node.Out, between)
	inst.WritePinBool(node.OutNot, !between)
	inst.WritePinBool(node.Above, above)
	inst.WritePinBool(node.Below, below)
}
