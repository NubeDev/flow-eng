package count

import (
	"github.com/NubeDev/flow-eng/node"
)

type CountString struct {
	*node.Spec
	count float64
}

func NewCountString(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, countStringNode, Category)
	cov := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, true)
	body.Inputs = node.BuildInputs(cov)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &CountString{body, 0}, nil
}

func (inst *CountString) Process() {
	updated, _, _ := inst.InputUpdated(node.In)
	if updated {
		inst.count++
	}
	inst.WritePinFloat(node.CountOut, inst.count)
}

func (inst *CountString) GetPersistedData() any {
	return &inst.count
}
