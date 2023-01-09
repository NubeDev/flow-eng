package statistics

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type MinMaxAvg struct {
	*node.Spec
}

func NewMinMaxAvg(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, minMaxAvg, category)
	if err != nil {
		return nil, err
	}
	min := node.BuildOutput(node.MinOutput, node.TypeFloat, nil, body.Outputs)
	max := node.BuildOutput(node.MaxOutput, node.TypeFloat, nil, body.Outputs)
	avg := node.BuildOutput(node.AvgOutput, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(min, max, avg)
	body = node.BuildNode(body, body.Inputs, outputs, nil)
	return &MinMaxAvg{body}, nil
}

func (inst *MinMaxAvg) Process() {
	count := inst.InputsLen()
	fmt.Println("MIN-MAX-AVG Process() input count: ", count)
	inputs := float.ConvertInterfaceToFloatMultiple(inst.ReadMultiple(count))
	var nonNilValues []float64
	for _, value := range inputs {
		if value != nil {
			nonNilValues = append(nonNilValues, *value)
		}
	}
	if len(nonNilValues) == 0 {
		inst.WritePinNull(node.MinOutput)
		inst.WritePinNull(node.MaxOutput)
		inst.WritePinNull(node.AvgOutput)
	} else {
		minValue := array.MinFloat64(nonNilValues)
		maxValue := array.MaxFloat64(nonNilValues)
		avgValue := average(nonNilValues)
		inst.WritePin(node.MinOutput, minValue)
		inst.WritePin(node.MaxOutput, maxValue)
		inst.WritePin(node.AvgOutput, avgValue)
	}
}
