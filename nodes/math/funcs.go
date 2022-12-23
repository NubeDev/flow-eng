package math

import "github.com/NubeDev/flow-eng/helpers/array"

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
	case add:
		output = array.Add(nonNilValues)
	case sub:
		output = array.Subtract(nonNilValues)
	case multiply:
		output = array.Multiply(nonNilValues)
	}
	return &output
}
