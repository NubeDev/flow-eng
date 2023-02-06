package numtransform

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type Limit struct {
	*node.Spec
}

func NewLimit(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, limitNode, category)
	in := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)
	min := node.BuildInput(node.MinInput, node.TypeFloat, nil, body.Inputs, nil)
	max := node.BuildInput(node.MaxInput, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in, min, max)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Limit{body}, nil
}

func (inst *Limit) Process() {
	in, _ := inst.ReadPinAsFloat(node.Inp)
	min, _ := inst.ReadPinAsFloat(node.MinInput)
	max, _ := inst.ReadPinAsFloat(node.MaxInput)
	inst.WritePinFloat(node.Out, float.LimitToRange(in, min, max))
}
