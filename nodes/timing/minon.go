package timing

import (
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type MinOn struct {
	*node.Spec
	timer     timer.TimedDelay
	lastValue bool
	locked    bool
}

func NewMinOn(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayMinON, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in, interval)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &MinOn{body, timer, false, false}, nil
}

func (inst *MinOn) run(interval int) {
	inst.locked = true
	time.Sleep(time.Duration(interval) * time.Second)
	inst.locked = false

}

func (inst *MinOn) Process() {
	in := inst.ReadPinBool(node.In)
	interval := inst.ReadPinAsInt(node.Interval)
	_, boolCOV := inst.InputUpdated(node.In)
	if boolCOV {
		if !inst.locked {
			go inst.run(interval)
		}
	}
	if inst.locked {
		inst.WritePin(node.Out, true)
	} else {
		inst.WritePin(node.Out, in)
	}

}

func (inst *MinOn) Cleanup() {}
