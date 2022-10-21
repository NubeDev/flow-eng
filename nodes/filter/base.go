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

	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return body, nil
}

func process(body node.Node) {
	//setting := body.GetSetting(node.Operation)
	//
	//setting.Properties

	in, _ := body.ReadPinAsFloat(node.In)
	output, ok := operation("ceil", in)
	if !ok {
		body.WritePin(node.Out, nil)
	} else {
		body.WritePin(node.Out, output)
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
