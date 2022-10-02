package bool

import (
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type MinOn struct {
	*node.Spec
	timer           timer.TimedDelay
	lastValue       bool
	triggered       bool
	lockedMinOn     bool
	lockedMinOff    bool
	minOffTriggered bool
}

func NewMinOn(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayMinOnOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	onInterval := node.BuildInput(node.OnInterval, node.TypeFloat, nil, body.Inputs)
	offInterval := node.BuildInput(node.OffInterval, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in, onInterval, offInterval)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &MinOn{body, timer, false, false, false, false, false}, nil
}

func (inst *MinOn) minOn(interval int) {
	inst.lockedMinOn = true
	time.Sleep(time.Duration(interval) * time.Second)
	inst.lockedMinOn = false
	inst.triggered = true

}

func (inst *MinOn) minOff(interval int) {
	inst.lockedMinOff = true
	time.Sleep(time.Duration(interval) * time.Second)
	inst.lockedMinOff = false
	inst.triggered = false

}

// out is true if in == true
// out is true if minOn is on
// trigger minOn if in == true and minOff is not active

func (inst *MinOn) Process() {
	in := inst.ReadPinBool(node.In)
	onInterval := inst.ReadPinAsInt(node.OnInterval)
	offInterval := inst.ReadPinAsInt(node.OffInterval)
	if in {
		if !inst.lockedMinOn && !inst.lockedMinOff && !inst.triggered {
			go inst.minOn(onInterval) // trigger minOn
		}
	}
	if inst.triggered && !in { // trigger minOff
		if !inst.lockedMinOff {
			go inst.minOff(offInterval)
		}

	}
	if inst.lockedMinOn {
		inst.WritePin(node.Out, true)
		return
	} else if inst.triggered && in { // minOn has finished and input is still true
		inst.WritePin(node.Out, true)
		return
	} else if inst.lockedMinOff {
		inst.WritePin(node.Out, false)
		return
	}
	if inst.triggered && !in {
		inst.triggered = false
		inst.WritePin(node.Out, false)
	} else {
		inst.WritePin(node.Out, false)
	}

}

func (inst *MinOn) Cleanup() {}
