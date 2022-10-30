package statistics

import (
	"github.com/NubeDev/flow-eng/node"
)

func average(xs []float64) float64 {
	var total float64
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}

type Avg struct {
	*node.Spec
}

func NewAvg(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, avg, category)
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 2, 2, 20, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Max{body}, nil
}

func (inst *Avg) Process() {
	Process(inst)

}
