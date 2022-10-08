package constant

import (
	"github.com/NubeDev/flow-eng/node"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, constNum, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Number{body}, nil
}

func (inst *Number) Process() {
	in1, null := inst.ReadPinAsFloatOk(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinFloat(node.Out, in1)
	}

}

func (inst *Number) Cleanup() {}
