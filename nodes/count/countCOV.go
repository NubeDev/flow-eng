package count

import (
	"github.com/NubeDev/flow-eng/node"
)

type CountCOV struct {
	*node.Spec
	count float64
}

func NewCountCOV(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, countCOVNode, category)
	cov := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(cov)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &CountCOV{body, 0}, nil
}

func (inst *CountCOV) Process() {
	updated, _ := inst.InputUpdated(node.CountUp)
	if updated {
		inst.count++
	}
	inst.WritePinFloat(node.CountOut, inst.count)

}
