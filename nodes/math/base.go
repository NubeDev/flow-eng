package math

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/mitchellh/mapstructure"
)

const (
	category = "math"
	divide   = "divide"
	add      = "add"
	sub      = "subtract"
	multiply = "multiply"
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
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, count, 2, 20, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(schemas.GetInputCount())
	body.SetDynamicInputs()
	return body, nil
}

func process(body node.Node) {
	equation := body.GetName()
	fmt.Println("MATH Process() node: ", equation)
	count := body.InputsLen()
	fmt.Println("MATH Process() input count: ", count)
	inputs := float.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	fmt.Println("MATH Process() inputs: ", inputs)
	for i, val := range inputs {
		if val != nil {
			fmt.Println("MATH Process(): ", i, " ", *val)
		} else {
			fmt.Println("MATH Process(): ", i, " ", val)
		}
	}
	output := operation(equation, inputs)
	if output == nil {
		body.WritePinNull(node.Out)
	} else {
		body.WritePinFloat(node.Out, float.NonNil(output))
	}
}
