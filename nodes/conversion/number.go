package conversion

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
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
	in1 := inst.ReadPinAsFloat(node.In)
	f, ok := conversions.GetFloatOk(in1)
	if ok {
		if f == 1 {
			inst.WritePin(node.Boolean, true)
		} else {
			inst.WritePin(node.Boolean, false)
		}
		inst.WritePin(node.String, fmt.Sprintf("%f", f))
	} else {
		inst.WritePin(node.Boolean, nil)
		inst.WritePin(node.String, nil)
	}
}

func (inst *Number) Cleanup() {}
