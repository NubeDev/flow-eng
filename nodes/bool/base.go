package bool

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

const (
	category = "bool"
)

const (
	and    = "and"
	or     = "or"
	not    = "not"
	xor    = "xor"
	toggle = "toggle"
)

const (
	inputCount = "Inputs Count"
)

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs) // TODO: this input shouldn't have a manual override value
	in2 := node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(in1, in2)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return body, nil
}

func Process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		body.WritePin(node.Out, nil)
	} else {
		// log.Infof("bool: %s, result: %v", equation, *output)
		body.WritePin(node.Out, float.NonNil(output))
	}
}

func operation(operation string, values []*float64) *float64 {
	var nonNilValues []float64
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
		if array.AllTrueFloat64(nonNilValues) {
			return float.New(1)
		} else {
			return float.New(0)
		}
	case or:
		if array.OneIsTrueFloat64(nonNilValues) {
			return float.New(1)
		} else {
			return float.New(0)
		}
	case not:
		if nonNilValues[0] == 0 {
			return float.New(1)
		} else {
			return float.New(0)
		}
	}
	return float.New(0)
}
