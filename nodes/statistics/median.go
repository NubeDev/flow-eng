package statistics

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"math"
	"sort"
)

type Median struct {
	*node.Spec
}

func NewMedian(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, median, category)
	if err != nil {
		return nil, err
	}
	return &Median{body}, nil
}

func (inst *Median) Process() {
	count := inst.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(inst.ReadMultiple(count))
	var nonNilValues []float64
	for _, value := range inputs {
		if value != nil {
			nonNilValues = append(nonNilValues, *value)
		}
	}
	if len(nonNilValues) == 0 {
		inst.WritePinNull(node.Out)
	} else if len(nonNilValues) == 1 {
		inst.WritePin(node.Out, nonNilValues[0])
	} else {
		sort.Float64s(nonNilValues)
		mid := float64(len(nonNilValues) / 2)
		if math.Mod(mid, 1) > 0 { // Odd number of elements in the input array, take the middle one
			inst.WritePin(node.Out, nonNilValues[int(mid-0.5)])
		} else { // Even number of elements in the input array, take the average of the middle 2 values
			averageOfMiddles := (nonNilValues[int(mid-float64(1))] + nonNilValues[int(mid)]) / 2
			inst.WritePin(node.Out, averageOfMiddles)
		}
	}
}
