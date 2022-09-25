package math

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/integer"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
)

const (
	constNum = "const-num"
	category = "math"
	divide   = "divide"
	add      = "add"
	sub      = "subtract"
	multiply = "multiply"
)

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	buildCount, setting, value, err := node.NewSetting(body, &node.SettingOptions{Type: node.Number, Title: node.InputCount, Min: 2, Max: 20})
	if err != nil {
		return nil, err
	}
	settings, err := node.BuildSettings(setting)
	if err != nil {
		return nil, err
	}
	count, ok := value.(int)
	if !ok {
		count = 2
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, count, integer.NonNil(buildCount.Min), integer.NonNil(buildCount.Max), body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Result, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	body.SetSchema(schemas.GetInputCount())
	return body, nil
}

func process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		//log.Infof("equation: %s, result: %v", equation, output)
	} else {
		//log.Infof("equation: %s, result: %v", equation, *output)
		body.WritePin(node.Result, float.NonNil(output))
	}

}
