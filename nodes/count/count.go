package count

import (
	"github.com/NubeDev/flow-eng/node"
)

type CountBool struct {
	*node.Spec
	count float64
}

func NewCountBool(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, countBoolNode, category)
	countUp := node.BuildInput(node.CountUp, node.TypeBool, nil, body.Inputs, nil)
	countDown := node.BuildInput(node.CountDown, node.TypeBool, nil, body.Inputs, nil)
	body.Inputs = node.BuildInputs(countUp, countDown)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &CountBool{body, 0}, nil
}

func (inst *CountBool) Process() {
	_, countUpUpdated := inst.InputUpdated(node.CountUp)
	if countUpUpdated {
		inst.count++
	}
	_, countDownUpdated := inst.InputUpdated(node.CountDown)
	if countDownUpdated {
		inst.count--
	}
	inst.WritePinFloat(node.CountOut, inst.count)

}
