package statistics

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"sort"
)

type Rank struct {
	*node.Spec
}

func NewRank(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, rank, category)
	if err != nil {
		return nil, err
	}
	out1 := node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs)
	out2 := node.BuildOutput(node.Out2, node.TypeFloat, nil, body.Outputs)
	out3 := node.BuildOutput(node.Out3, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out1, out2, out3)
	body = node.BuildNode(body, body.Inputs, outputs, nil)
	return &Rank{body}, nil
}

func (inst *Rank) Process() {
	count := inst.InputsLen()
	inputs := float.ConvertInterfaceToFloatMultiple(inst.ReadMultiple(count))
	var nonNilValues []float64
	for _, value := range inputs {
		if value != nil {
			nonNilValues = append(nonNilValues, *value)
		}
	}
	if len(nonNilValues) == 0 {
		inst.WritePinNull(node.Out1)
		inst.WritePinNull(node.Out2)
		inst.WritePinNull(node.Out3)
	} else {
		sort.Float64s(nonNilValues)
		// TODO: Add in minimum rank
		if len(nonNilValues) == 1 {
			inst.WritePin(node.Out1, nonNilValues[0])
			inst.WritePinNull(node.Out2)
			inst.WritePinNull(node.Out3)
		} else if len(nonNilValues) == 2 {
			inst.WritePin(node.Out1, nonNilValues[1])
			inst.WritePin(node.Out2, nonNilValues[0])
			inst.WritePinNull(node.Out3)
		} else if len(nonNilValues) >= 3 {
			inst.WritePin(node.Out1, nonNilValues[len(nonNilValues)-1])
			inst.WritePin(node.Out2, nonNilValues[len(nonNilValues)-2])
			inst.WritePin(node.Out3, nonNilValues[len(nonNilValues)-3])
		}
	}
}
