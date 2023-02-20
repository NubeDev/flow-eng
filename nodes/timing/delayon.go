package timing

import (
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type DelayOn struct {
	*node.Spec
	timer      *time.Timer
	currOutput bool
	lastDelay  time.Duration
}

func NewDelayOn(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, delayOn, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false, false) // TODO: this input shouldn't have a manual override value
	delayInput := node.BuildInput(node.Interval, node.TypeFloat, 1, body.Inputs, true, false)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false, false) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(in, delayInput, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildDefaultSchema())
	return &DelayOn{body, nil, false, 0}, nil
}

/*
if node input is true
start delay, after the delay set the output to true
*/

func (inst *DelayOn) Process() {
	delayDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if delayDuration != inst.lastDelay {
		inst.setSubtitle(delayDuration)
		inst.lastDelay = delayDuration
	}

	in1, _ := inst.ReadPinAsBool(node.In)
	if !in1 { // any time input is false, set output false and cancel any running timers
		inst.WritePinFalse(node.Out)
		inst.currOutput = false
		if inst.timer != nil {
			inst.timer.Stop()
			inst.timer = nil
		}
		return
	}

	// input is true

	if inst.currOutput { // input is still active, so output is still active, cancel any running timers (for safety)
		inst.WritePinTrue(node.Out)
		inst.currOutput = true
		if inst.timer != nil {
			inst.timer.Stop()
			inst.timer = nil
		}
		return
	}

	// input is active, but output isn't so start a timer if it doesn't exist already
	if inst.timer == nil {
		inst.timer = time.AfterFunc(delayDuration, func() {
			inst.WritePinTrue(node.Out)
			inst.currOutput = true
			inst.timer = nil
		})
	}
	inst.WritePinBool(node.Out, inst.currOutput)
}

func (inst *DelayOn) Start() {
	inst.WritePinFalse(node.Out)
	inst.currOutput = false
}

func (inst *DelayOn) Stop() {
	if inst.timer != nil {
		inst.timer.Stop()
		inst.timer = nil
	}
}

func (inst *DelayOn) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}
