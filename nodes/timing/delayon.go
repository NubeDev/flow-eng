package timing

import (
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type DelayOn struct {
	*node.Spec
	timer     timer.TimedDelay
	triggered bool
	active    bool
}

func NewDelayOn(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOn, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	delayTime := node.BuildInput(node.DelaySeconds, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(delayTime, in)
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	return &DelayOn{body, timer, false, false}, nil
}

/*
if node input is true
start delay, after the delay set the triggered to true
*/

func duration(f int) time.Duration {
	return time.Duration(f) * time.Second
}

func (inst *DelayOn) Process() {
	timeDelay, _ := inst.ReadPinAsInt(node.DelaySeconds)

	in1, _ := inst.ReadPinAsFloat(node.In)

	if in1 >= 1 && inst.active { // timer has gone to true and input is still true
		inst.WritePin(node.Out, 1)
		return
	} else {
		inst.active = false // timer is true and input went back to 0
	}

	if in1 < 1 && inst.triggered { // went true but not for long enough to finish the timeOn delay, so reset the timer
		inst.timer.Stop()
		inst.timer = timer.NewTimer()
		inst.triggered = false
	}
	// fmt.Println(timeDelay, time.Duration(timeDelay)*time.Second, 99999, time.Duration(timeDelay).Seconds())

	if in1 >= 1 {
		if !inst.timer.WaitFor(duration(timeDelay)) {
			inst.WritePin(node.Out, 0)
			inst.triggered = true
			return
		} else {
			inst.active = true
			inst.WritePin(node.Out, 1)
		}

	} else {
		inst.WritePin(node.Out, 0)
	}

}

func (inst *DelayOn) Cleanup() {}
