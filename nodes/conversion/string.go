package conversion

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"strconv"
)

type String struct {
	*node.Spec
}

func NewString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, conversionString, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs, nil))
	asBool := node.BuildOutput(node.Boolean, node.TypeBool, nil, body.Outputs)
	asString := node.BuildOutput(node.Float, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(asBool, asString)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &String{body}, nil
}

func (inst *String) Process() {
	in1, null := inst.ReadPinAsString(node.In)
	if null {
		inst.WritePinNull(node.Float)
		inst.WritePinNull(node.Boolean)
		return
	}
	f, ok := conversions.GetFloatOk(in1)
	if ok { // to float
		inst.WritePinFloat(node.Float, f)
	} else {
		inst.WritePinNull(node.Float)
	}
	result, err := strconv.ParseBool(in1) // to bool
	if err != nil {
		inst.WritePinNull(node.Boolean)
	} else {
		inst.WritePinBool(node.Boolean, result)
	}

}
