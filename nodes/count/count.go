package count

import (
	"github.com/NubeDev/flow-eng/node"
)

type Count struct {
	*node.Spec
	count uint64
}

func NewCount(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, countNode, category)
	countUp := node.BuildInput(node.CountUp, node.TypeBool, nil, body.Inputs)
	countDown := node.BuildInput(node.CountDown, node.TypeBool, nil, body.Inputs)
	body.Inputs = node.BuildInputs(countUp, countDown)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &Count{body, 0}, nil
}

func (inst *Count) Process() {
	_, countUpUpdated := inst.InputUpdated(node.CountUp)
	if countUpUpdated {
		inst.count = inst.count + 1
	}
	_, countDownUpdated := inst.InputUpdated(node.CountDown)
	if countDownUpdated {
		inst.count = inst.count - 1
	}

	inst.WritePin(node.CountOut, inst.count)

}
