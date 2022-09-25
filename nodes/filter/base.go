package filter

import (
	"github.com/NubeDev/flow-eng/node"
	"math"
)

const (
	category = "math-advanced"
	sin      = "sin"
)

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)

	_, setting, _, err := node.NewSetting(body, &node.SettingOptions{Type: node.Array, Title: node.Operation})
	if err != nil {
		return nil, err
	}
	settings, err := node.BuildSettings(setting)
	if err != nil {
		return nil, err
	}
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return body, nil
}

func process(body node.Node) {
	//setting := body.GetSetting(node.Operation)
	//
	//setting.Properties

	in := body.ReadPinAsFloat(node.In)
	output, ok := operation("ceil", in)
	if !ok {
		body.WritePin(node.Result, nil)
	} else {
		body.WritePin(node.Result, output)
	}
}

func operation(operation string, value float64) (val float64, ok bool) {
	output := 0.0
	switch operation {
	case "ceil":
		output = math.Ceil(value)
		ok = true
	}
	return output, ok
}
