package math

import (
	"github.com/NubeDev/flow-eng/node"
	"math"
)

type Modulo struct {
	*node.Spec
}

func NewModulo(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, modulo, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs, nil)
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in1, in2)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &Modulo{body}, nil
}

func (inst *Modulo) Process() {
	in1, null1 := inst.ReadPinAsFloat(node.In1)
	in2, null2 := inst.ReadPinAsFloat(node.In2)
	if null1 || null2 || in2 == 0 {
		inst.WritePinNull(node.Out)
	} else {
		output := math.Mod(in1, in2)
		inst.WritePinFloat(node.Out, output)
	}

}
