package timing

import (
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"log"
	"time"
)

type DelayOn struct {
	*node.Spec
	timer   flowctrl.TimedDelay
	waiting bool
	active  bool
}

func NewDelayOn(body *node.Spec, timer flowctrl.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOn, category)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	return &DelayOn{body, timer, false, false}, nil
}

/*
if node input is true
start delay, after the delay set the triggered to true
*/

func (inst *DelayOn) Process() {
	log.Println("Delayed START")
	in1 := inst.ReadPinAsFloat(node.In)

	if in1 >= 1 && inst.active { // timer has gone to true and input is still true
		inst.WritePin(node.Out, 1)
		return
	} else {
		inst.active = false // timer is true and input went back to 0
	}

	if in1 >= 1 {
		if !inst.timer.WaitFor(5 * time.Second) {
			inst.WritePin(node.Out, 0)
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
