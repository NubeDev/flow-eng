package conversion

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"strconv"
)

type String struct {
	*node.Spec
}

func NewString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, conversionString, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs))
	asBool := node.BuildOutput(node.Boolean, node.TypeBool, nil, body.Outputs)
	asString := node.BuildOutput(node.Float, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(asBool, asString)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &String{body}, nil
}

func (inst *String) Process() {
	in1, _ := inst.ReadPinAsString(node.In)
	f, ok := conversions.GetFloatOk(in1)
	if ok { // to float
		inst.WritePin(node.Float, fmt.Sprintf("%f", f))
	} else {
		inst.WritePin(node.Float, nil)
	}

	result, err := strconv.ParseBool(in1) // to bool
	if err != nil {
		inst.WritePin(node.Boolean, nil)
	} else {
		inst.WritePin(node.Boolean, result)
	}

}
