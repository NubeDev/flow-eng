package trigger

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type Random struct {
	*node.Spec
}

func NewRandom(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, randomFloat, category)
	min := node.BuildInput(node.Min, node.TypeFloat, nil, body.Inputs)
	max := node.BuildInput(node.Max, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(min, max)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &Random{body}, nil
}

func (inst *Random) Process() {
	min := inst.ReadPinAsFloat(node.Min)
	max := inst.ReadPinAsFloat(node.Max)
	inst.WritePin(node.Out, float.RandFloat(min, max))

}
func (inst *Random) Cleanup() {}
