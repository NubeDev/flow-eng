package trigger

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type Random struct {
	*node.Spec
	precision int
}

func NewRandom(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, randomFloat, category)
	min := node.BuildInput(node.Min, node.TypeFloat, nil, body.Inputs)
	max := node.BuildInput(node.Max, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(min, max)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	precision, _ := getSettings(body.GetSettings())
	return &Random{body, precision}, nil
}

func (inst *Random) Process() {
	min, _ := inst.ReadPinAsFloat(node.Min)
	max, _ := inst.ReadPinAsFloat(node.Max)
	f := float.RandFloat(min, max)
	inst.WritePinFloat(node.Out, f, inst.precision)
}
