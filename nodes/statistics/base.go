package statistics

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/mitchellh/mapstructure"
)

const (
	category = "statistics"

	max       = "max"
	min       = "min"
	avg       = "average"
	minMaxAvg = "min-max-avg"
	rangeNode = "range"
	rank      = "rank"
	median    = "median"
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
	count := body.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		body.WritePinNull(node.Out)
	} else {
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
	output := 0.0
	switch operation {
	case min:
		output = array.MinFloat64(nonNilValues)
	case max:
		output = array.MaxFloat64(nonNilValues)
	case avg:
		output = average(nonNilValues)
	}
	return &output
}

func average(inputArray []float64) float64 {
	var total float64
	for _, v := range inputArray {
		total += v
	}
	return total / float64(len(inputArray))
}
