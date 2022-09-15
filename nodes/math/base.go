package math

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

const (
	constNum = "const-num"
	category = "math"
	divide   = "divide"
	add      = "add"
	sub      = "subtract"
	multiply = "multiply"
)

func Process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		log.Infof("equation: %s, result: %v", equation, output)
	} else {
		log.Infof("equation: %s, result: %v", equation, *output)
	}
	body.WritePin(node.Result, output)
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
	case add:
		output = array.Add(nonNilValues)
	case sub:
		output = array.Subtract(nonNilValues)
	case multiply:
		output = array.Multiply(nonNilValues)
	case divide:
		output = array.Divide(nonNilValues)
	}
	return &output
}
