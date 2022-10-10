package conversion

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, conversionNum, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	asBool := node.BuildOutput(node.Boolean, node.TypeBool, nil, body.Outputs)
	asString := node.BuildOutput(node.String, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(asBool, asString)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Number{body}, nil
}

func (inst *Number) Process() {
	in1, null := inst.ReadPinAsFloat(node.In)
	if !null {
		if in1 == 1 {
			inst.WritePinBool(node.Boolean, true)
		} else {
			inst.WritePinBool(node.Boolean, false)
		}
		inst.WritePin(node.String, fmt.Sprintf("%f", in1))
	} else {
		inst.WritePinNull(node.Boolean)
		inst.WritePinNull(node.String)
	}
}
