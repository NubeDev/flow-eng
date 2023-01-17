package timing

import (
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type DelayOff struct {
	*node.Spec
	timer      *time.Timer
	currOutput bool
	lastDelay  time.Duration
}

func NewDelayOff(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOff, category)
	in := node.BuildInput(node.Inp, node.TypeBool, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	delayInput := node.BuildInput(node.Delay, node.TypeFloat, nil, body.Inputs, str.New("interval"))
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(in, delayInput, reset)

	out := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	body.SetSchema(buildDefaultSchema())

	return &DelayOff{body, nil, false, 1 * time.Second}, nil
}

func (inst *DelayOff) Process() {
	delayDuration, _ := inst.ReadPinAsTimeSettings(node.Delay)
	if delayDuration != inst.lastDelay {
		inst.setSubtitle(delayDuration)
		inst.lastDelay = delayDuration
	}

	in1, _ := inst.ReadPinAsBool(node.In)
	if in1 { // any time input is true, set output true and cancel any running timers
		inst.WritePinTrue(node.Outp)
		inst.currOutput = true
		if inst.timer != nil {
			inst.timer.Stop()
			inst.timer = nil
		}
		return
	}

	// input is false

	if !inst.currOutput { // input is still false, so output is still false, cancel any running timers (for safety)
		inst.WritePinFalse(node.Outp)
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
			inst.WritePinFalse(node.Outp)
			inst.currOutput = false
			inst.timer = nil
		})
	}
}

func (inst *DelayOff) Start() {
	inst.WritePinFalse(node.Outp)
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
