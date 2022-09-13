package math

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

const (
	category = "math"
)

const (
	and     = "and"
	or      = "or"
	not     = "not"
	greater = "greater"
	less    = "less"
)

const (
	inputCount = "Inputs Count"
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
	body.WritePin(node.Out1, output)
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
