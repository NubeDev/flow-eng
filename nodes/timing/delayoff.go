package timing

import (
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type DelayOff struct {
	*node.Spec
	timer      *time.Timer
	currOutput bool
	lastDelay  time.Duration
}

func NewDelayOff(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, delayOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false) // TODO: this input shouldn't have a manual override value
	delayInput := node.BuildInput(node.Interval, node.TypeFloat, 1, body.Inputs, true)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(in, delayInput, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildDefaultSchema())

	return &DelayOff{body, nil, false, 0}, nil
}

func (inst *DelayOff) Process() {
	delayDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if delayDuration != inst.lastDelay {
		inst.setSubtitle(delayDuration)
		inst.lastDelay = delayDuration
	}

	in1, _ := inst.ReadPinAsBool(node.In)
	if in1 { // any time input is true, set output true and cancel any running timers
		inst.WritePinTrue(node.Out)
		inst.currOutput = true
		if inst.timer != nil {
			inst.timer.Stop()
			inst.timer = nil
		}
		return
	}

	// input is false

	if !inst.currOutput { // input is still false, so output is still false, cancel any running timers (for safety)
		inst.WritePinFalse(node.Out)
		inst.currOutput = false
		if inst.timer != nil {
			inst.timer.Stop()
			inst.timer = nil
		}
		return
	}

	// input is false, but output isn't so start a timer if it doesn't exist already
	if inst.timer == nil {
		inst.timer = time.AfterFunc(delayDuration, func() {
			inst.WritePinFalse(node.Out)
			inst.currOutput = false
			inst.timer = nil
		})
	}
	inst.WritePinBool(node.Out, inst.currOutput)
}

func (inst *DelayOff) Start() {
	inst.WritePinFalse(node.Out)
	inst.currOutput = false
}

func (inst *DelayOff) Stop() {
	if inst.timer != nil {
		inst.timer.Stop()
		inst.timer = nil
	}
}

func (inst *DelayOff) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}
