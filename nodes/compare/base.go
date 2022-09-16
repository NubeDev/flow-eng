package compare

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

const (
	category = "compare"
)

const (
	logicCompare = "compare"
	between      = "between"
)

func B2F(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func zeroToOne(b float64) float64 {
	if b > 0 {
		return 0
	}
	return 1
}

func Process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	val1, val2, val3, val4 := operation(equation, inputs)
	if equation == logicCompare {
		body.WritePin(node.GraterThan, val1)
		body.WritePin(node.LessThan, val2)
		body.WritePin(node.Equal, val3)
		if val1 == nil {
			log.Infof("compare: %s, result-%s: %v", equation, logicCompare, val1)
		} else {
			log.Infof("compare: %s, result-%s: %v", equation, logicCompare, *val1)
		}
	}
	if equation == between {
		body.WritePin(node.Out, val1)
		body.WritePin(node.OutNot, val2)
		body.WritePin(node.Above, val3)
		body.WritePin(node.Below, val4)
		if val1 == nil {
			log.Infof("compare: %s, result-%s: %v", equation, between, val1)
		} else {
			log.Infof("compare: %s, result-%s: %v", equation, between, *val1)
		}
	}

}

func operation(operation string, values []*float64) (*float64, *float64, *float64, *float64) {
	var nonNilValues []float64
	for _, value := range values {
		if value != nil {
			nonNilValues = append(nonNilValues, *value)
		}
	}
	if len(nonNilValues) == 0 {
		return nil, nil, nil, nil
	}
	switch operation {
	case logicCompare:
		greater, less, equal := array.Compare(nonNilValues)
		return float.New(B2F(greater)), float.New(B2F(less)), float.New(B2F(equal)), nil
	case between:
		if len(nonNilValues) == 3 {
			between, below, above := array.Between(nonNilValues[0], nonNilValues[1], nonNilValues[2])
			outNot := float.New(zeroToOne(B2F(above)))
			return float.New(B2F(between)), outNot, float.New(B2F(below)), float.New(B2F(above))
		}
	}
	return nil, nil, nil, nil
}
