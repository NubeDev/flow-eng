package mathematics

import (
	"github.com/NubeDev/flow-eng/node"
)

const (
	category     = "math"
	mathAdvanced = "advanced"
)

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return body, nil
}

func process(body node.Node) {
	function, err := getSettings(body.GetSettings())
	if err != nil {
		return
	}
	if function == "" {
		function = acos
	}
	in, null := body.ReadPinAsFloat(node.In)
	if null {
		body.WritePinNull(node.Result)
	}
	output, err := mathFunc(function, in)
	if err != nil {
		body.WritePin(node.Result, 0)
	} else {
		body.WritePin(node.Result, output)
	}
}
