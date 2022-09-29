package timing

import (
	timer "github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type Delay struct {
	*node.Spec
	timer timer.TimedDelay
}

func NewDelay(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delay, category)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	return &Delay{body, timer}, nil
}

func (inst *Delay) Process() {
	in1 := inst.ReadPinAsFloat(node.In)
	if !inst.timer.WaitFor(5 * time.Second) {
		inst.WritePin(node.Out, 0)
		return
	}
	inst.WritePin(node.Out, in1)

}

func (inst *Delay) Cleanup() {}
