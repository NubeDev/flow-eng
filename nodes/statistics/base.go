package statistics

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

const (
	category = "statistics"

	max = "min"
	min = "max"
	avg = "avg"
)

func Process(body node.Node) {
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
