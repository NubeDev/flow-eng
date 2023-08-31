package conversion

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
)

type Number struct {
	*node.Spec
}

func NewNumber(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, conversionNum, Category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false))
	asBool := node.BuildOutput(node.Boolean, node.TypeBool, nil, body.Outputs)
	asString := node.BuildOutput(node.String, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(asBool, asString)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Number{body}, nil
}

func (inst *Number) Process() {
	in1, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Boolean)
		inst.WritePinNull(node.String)
	} else {
		if in1 == 1 {
			inst.WritePinBool(node.Boolean, true)
		} else {
			inst.WritePinBool(node.Boolean, false)
		}
		v := conversions.FloatToString(in1)
		if v != "" {
			inst.WritePin(node.String, conversions.FloatToString(in1))
		} else {
			inst.WritePinNull(node.String)
		}
	}
}

type Str2Num struct {
	*node.Spec
}

func NewStr2Num(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, conversionString2Num, Category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false))
	asFloat := node.BuildOutput(node.Float, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(asFloat)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Str2Num{body}, nil
}

func (inst *Str2Num) Process() {
	in1, null := inst.ReadPinAsString(node.In)
	if null {
		inst.WritePinNull(node.Float)
		return
	}
	f, ok := conversions.GetFloatOk(in1)
	if ok { // to float
		inst.WritePinFloat(node.Float, f)
	} else {
		inst.WritePinNull(node.Float)
	}
}

type Num2Str struct {
	*node.Spec
}

func NewNum2Str(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, conversionNum2Str, Category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false))
	asString := node.BuildOutput(node.String, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(asString)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Num2Str{body}, nil
}

func (inst *Num2Str) Process() {
	in1, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.String)
	} else {
		v := conversions.FloatToString(in1)
		if v != "" {
			inst.WritePin(node.String, conversions.FloatToString(in1))
		} else {
			inst.WritePinNull(node.String)
		}
	}
}
