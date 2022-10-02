package mathematics

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/mitchellh/mapstructure"
)

const (
	category     = "math"
	mathAdvanced = "advanced"
)

type nodeSettings struct {
	Function string `json:"function"`
}

func getSettings(body map[string]interface{}) (string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return "", err
	}
	if settings != nil {
		return settings.Function, nil
	}
	return "", nil
}

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
	in := body.ReadPinAsFloat(node.In)
	output, err := mathFunc(function, in)
	if err != nil {
		body.WritePin(node.Result, 0)
	} else {
		body.WritePin(node.Result, output)
	}
}
