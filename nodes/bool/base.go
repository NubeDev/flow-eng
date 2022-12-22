package bool

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/mitchellh/mapstructure"
)

const (
	category = "bool"
)

const (
	and           = "and"
	or            = "or"
	not           = "not"
	xor           = "xor"
	toggle        = "toggle"
	delayMinOnOff = "min on off"
)

const (
	inputCount = "Inputs Count"
)

type nodeSettings struct {
	InputCount int `json:"inputCount"`
}

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	settings := &nodeSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	var count = 2
	if settings != nil {
		count = settings.InputCount
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeBool, nil, count, 2, 20, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(schemas.GetInputCount())
	body.SetDynamicInputs()
	return body, nil
}

func Process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := boolean.ConvertInterfaceToBoolMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		body.WritePinNull(node.Out)
	} else {
		body.WritePinBool(node.Out, boolean.NonNil(output))
	}
}

func operation(operation string, values []*bool) *bool {
	var nonNilValues []bool
	for _, value := range values {
		if value != nil {
			nonNilValues = append(nonNilValues, *value)
		}
	}
	if len(nonNilValues) == 0 {
		return nil
	}
	switch operation {
	case and:
		if array.AllTrue(nonNilValues) {
			return boolean.New(true)
		} else {
			return boolean.New(false)
		}
	case or:
		if array.OneIsTrue(nonNilValues) {
			return boolean.New(true)
		} else {
			return boolean.New(false)
		}
	case not:
		if nonNilValues[0] {
			return boolean.New(false)
		} else {
			return boolean.New(true)
		}
	}
	return boolean.New(false)
}
