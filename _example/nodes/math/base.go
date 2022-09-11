package math

import (
	"errors"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

const (
	constNum = "const-num"
	category = "math"
	add      = "add"
	sub      = "subtract"
	multiply = "multiply"
)

const (
	inputCount = "Inputs Count"
)

func inputsCount(body *node.BaseNode) (*node.PropertyBase, *node.Settings, int, error) {
	const min = 2
	const max = 20
	var count = min
	count = body.GetPropValueInt(inputCount, min)
	base := &node.PropertyBase{
		Min: min,
		Max: max,
	}
	setting, err := node.NewSetting(node.Number, inputCount, base)
	if err != nil {
		return nil, nil, 0, err
	}
	return base, setting, count, err
}

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
	case add:
		return array.Add(values), nil
	case sub:
		return array.Subtract(values), nil
	case multiply:
		return array.Multiply(values), nil
	}
	return 0, errors.New("invalid operation")
}
