package count

import (
	"github.com/NubeDev/flow-eng/node"
)

type CountNum struct {
	*node.Spec
	count float64
}

func NewCountNum(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, countNumNode, category)
	cov := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, true)
	body.Inputs = node.BuildInputs(cov)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &CountNum{body, 0}, nil
}

func (inst *CountNum) Process() {
	updated, _ := inst.InputUpdated(node.In)
	if updated {
		inst.count++
	}
	inst.WritePinFloat(node.CountOut, inst.count)

}
