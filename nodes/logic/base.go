package math

import (
	"errors"
	"github.com/NubeDev/flow-eng/helpers/array"
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
	f, err := operation(equation, body.ReadMultipleNums(count))
	if err != nil {
		body.WritePin(node.Out1, nil)
		return
	}
	log.Infof("equation:%s result:%f", equation, f)
	body.WritePinNum(node.Out1, f)
}

func operation(operation string, values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, errors.New("no values where pass in")
	}
	switch operation {
	case and:
		if array.AllTrueFloat64(values) {
			return 1, nil
		} else {
			return 0, nil
		}
	case or:
		if array.OneIsTrueFloat64(values) {
			return 1, nil
		} else {
			return 0, nil
		}
	case not:
		if len(values) > 0 {
			if values[0] == 0 {
				return 1, nil
			}
			if values[0] == 0 {
				return 1, nil
			} else {
				return 0, nil
			}
		}
	}
	return 0, errors.New("invalid operation")
}
