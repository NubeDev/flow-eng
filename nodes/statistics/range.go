package statistics

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
)

type Range struct {
	*node.Spec
}

func NewRange(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, rangeNode, Category)
	if err != nil {
		return nil, err
	}
	return &Range{body}, nil
}

func (inst *Range) Process() {
	count := inst.InputsLen()
	inputs := conversions.ConvertInterfaceToFloatMultiple(inst.ReadMultiple(count))
	var nonNilValues []float64
	for _, value := range inputs {
		if value != nil {
			nonNilValues = append(nonNilValues, *value)
		}
	}
	if len(nonNilValues) == 0 {
		inst.WritePinNull(node.Out)
	} else {
		minValue := array.MinFloat64(nonNilValues)
		maxValue := array.MaxFloat64(nonNilValues)
		rangeValue := maxValue - minValue
		inst.WritePinFloat(node.Out, rangeValue)
	}
}
