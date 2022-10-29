package timing

import (
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
)

type DelayOff struct {
	*node.Spec
	timer     timer.TimedDelay
	triggered bool
	active    bool
}

// TODO code is copied from delayOn, so needs to be finished

func NewDelayOff(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	delayTime := node.BuildInput(node.DelaySeconds, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(delayTime, in)
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	return &DelayOn{body, timer, false, false}, nil
}

func (inst *DelayOff) Process() {
	settings, _ := getSettings(inst.GetSettings())
	in1, null := inst.ReadPinAsBool(node.In)
	if null {
		inst.WritePinNull(node.Out)
	}

	if in1 && inst.active { // timer has gone to true and input is still true
		inst.WritePinTrue(node.Out)
		return
	} else {
		inst.active = false // timer is true and input went back to 0
	}

	if !in1 && inst.triggered { // went true but not for long enough to finish the timeOn delay, so reset the timer
		inst.timer.Stop()
		inst.timer = timer.NewTimer()
		inst.triggered = false
	}
	if in1 {
		if !inst.timer.WaitFor(duration(settings.Duration, settings.Time)) {
			inst.WritePinFalse(node.Out)
			inst.triggered = true
			return
		} else {
			inst.active = true
			inst.WritePinTrue(node.Out)
		}

	} else {
		inst.WritePinFalse(node.Out)
	}

}
