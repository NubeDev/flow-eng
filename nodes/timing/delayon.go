package timing

import (
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
)

type DelayOn struct {
	*node.Spec
	timer     timer.TimedDelay
	triggered bool
	active    bool
}

func NewDelayOn(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOn, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(in)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &DelayOn{body, timer, false, false}, nil
}

/*
if node input is true
start delay, after the delay set the triggered to true
*/

func (inst *DelayOn) Process() {
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
