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
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
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
	in, null := body.ReadPinAsFloat(node.Inp)
	if null {
		body.WritePinNull(node.Outp)
	}
	output, err := mathFunc(function, in)
	if err != nil {
		body.WritePinFloat(node.Outp, 0)
	} else {
		body.WritePinFloat(node.Outp, output)
	}
}
