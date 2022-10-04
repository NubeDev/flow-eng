package point

import (
	"github.com/NubeDev/flow-eng/node"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pointNumber, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in1, in2)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &Number{body}, nil
}

func (inst *Number) Process() {
	in1 := inst.ReadPinAsFloatPointer(node.In1)
	in2 := inst.ReadPinAsFloatPointer(node.In2)
	if in1 != nil {
		inst.WritePin(node.Out, in1)
	} else {
		inst.WritePin(node.Out, in2)
	}

}
func (inst *Number) Cleanup() {}
