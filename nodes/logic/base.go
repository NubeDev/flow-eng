package logic

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/integer"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

const (
	category = "logic"
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

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	buildCount, setting, value, err := node.NewSetting(body, &node.SettingOptions{Type: node.Number, Title: node.InputCount, Min: 2, Max: 20})
	if err != nil {
		return nil, err
	}
	settings, err := node.BuildSettings(setting)
	if err != nil {
		return nil, err
	}
	count, ok := value.(int)
	if !ok {
		count = 2
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, count, integer.NonNil(buildCount.Min), integer.NonNil(buildCount.Max), body.Inputs, node.ABCs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Result, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return body, nil
}

func Process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		log.Infof("logic: %s, result: %v", equation, output)
	} else {
		log.Infof("logic: %s, result: %v", equation, *output)
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
