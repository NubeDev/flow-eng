package conversion

import (
	"github.com/NubeDev/flow-eng/node"
)

type Boolean struct {
	*node.Spec
}

func NewBoolean(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, conversionBool, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeBool, nil, body.Inputs, false))
	asFloat := node.BuildOutput(node.Float, node.TypeFloat, nil, body.Outputs)
	asString := node.BuildOutput(node.String, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(asFloat, asString)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Boolean{body}, nil
}

func (inst *Boolean) Process() {
	in1, _ := inst.ReadPinAsBool(node.Inp)
	if in1 {
		inst.WritePinFloat(node.Float, 1)
		inst.WritePin(node.String, "1")
	} else {
		inst.WritePinFloat(node.Float, 0)
		inst.WritePin(node.String, "0")
	}
}
