package bool

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/NubeDev/flow-eng/node"
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
	delayMinOnOff = "min-on off"
)

const (
	inputCount = "Inputs Count"
)

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	in1 := node.BuildInput(node.In1, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value
	in2 := node.BuildInput(node.In2, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(in1, in2)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return body, nil
}

func Process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := boolean.ConvertInterfaceToBoolMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		body.WritePin(node.Out, nil)
	} else {
		// log.Infof("bool: %s, result: %v", equation, *output)
		body.WritePin(node.Out, boolean.NewFalse())
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
