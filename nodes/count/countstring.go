package count

import (
	"github.com/NubeDev/flow-eng/node"
)

type CountString struct {
	*node.Spec
	count float64
}

func NewCountString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, countStringNode, category)
	cov := node.BuildInput(node.Inp, node.TypeString, nil, body.Inputs, nil)
	body.Inputs = node.BuildInputs(cov)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &CountString{body, 0}, nil
}

func (inst *CountString) Process() {
	updated, _ := inst.InputUpdated(node.Inp)
	if updated {
		inst.count++
	}
	inst.WritePinFloat(node.CountOut, inst.count)
}
