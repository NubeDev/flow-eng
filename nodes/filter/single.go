package filter

import (
	"github.com/NubeDev/flow-eng/node"
)

type OnlyTrue struct {
	*node.Spec
}

func NewOnlyTrue(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyTrue, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &OnlyTrue{body}, nil
}

func (inst *OnlyTrue) Process() {
	v, _ := inst.ReadPinAsBool(node.In)
	if v {
		inst.WritePinTrue(node.Out)
	} else {
		inst.WritePinFalse(node.Out)
	}

}
