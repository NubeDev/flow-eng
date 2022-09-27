package switches

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/integer"
	"github.com/NubeDev/flow-eng/node"
)

const (
	selectNum  = "select-numeric"
	switchNode = "switch"
	category   = "switch"
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
	var nodeInputs []*node.Input
	selection := node.BuildInput(node.Selection, node.TypeFloat, nil, body.Inputs) // TODO: this input shouldn't have a manual override value
	inputsCount := node.DynamicInputs(node.TypeFloat, nil, count, integer.NonNil(buildCount.Min), integer.NonNil(buildCount.Max), body.Inputs, node.ABCs)
	nodeInputs = append(nodeInputs, selection)
	nodeInputs = append(nodeInputs, inputsCount...)
	inputs := node.BuildInputs(nodeInputs...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Result, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return body, nil
}

func process(body node.Node) {
	op := body.GetName()
	count := body.InputsLen()
	inputs := body.ReadMultiple(count)
	selection := body.ReadPinAsInt(node.Selection)
	if op == selectNum {
		output := selectValue(selection, inputs)
		if output == nil {
			body.WritePin(node.Result, nil)
		} else {
			body.WritePin(node.Result, conversions.GetFloat(output))
		}

	}

}

func selectValue(num int, values []interface{}) interface{} {
	if num == 0 {
		return nil
	}
	for i, value := range values {
		if num == i {
			return value
		}
	}
	return nil
}
