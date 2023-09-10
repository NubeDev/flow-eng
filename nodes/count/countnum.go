package count

import (
	"github.com/NubeDev/flow-eng/node"
)

type CountNum struct {
	*node.Spec
	count float64
}

func NewCountNum(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, countNumNode, Category)
	cov := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, true)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false, true)
	body.Inputs = node.BuildInputs(cov, reset)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &CountNum{body, 0}, nil
}

func (inst *CountNum) Process() {
	reset, _ := inst.ReadPinAsBool(node.Reset)
	if reset {
		inst.count = 0
	}
	updated, _, _ := inst.InputUpdated(node.In)
	if updated {
		inst.count++
	}
	inst.WritePinFloat(node.CountOut, inst.count)

}

func (inst *CountNum) GetPersistedData() any {
	return &inst.count
}
