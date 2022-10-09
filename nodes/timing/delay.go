package timing

import (
	timer "github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type Delay struct {
	*node.Spec
	timer     timer.TimedDelay
	lastValue float64
}

func NewDelay(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delay, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in, interval)
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	return &Delay{body, timer, 0}, nil
}

func (inst *Delay) Process() {
	in, _ := inst.ReadPinAsFloat(node.In)
	interval := inst.ReadPinAsInt(node.Interval)
	if !inst.timer.WaitFor(time.Duration(interval) * time.Second) {
		inst.WritePin(node.Out, inst.lastValue)
	} else {
		inst.lastValue = in
		inst.WritePin(node.Out, in)
	}

}

func (inst *Delay) Cleanup() {}
