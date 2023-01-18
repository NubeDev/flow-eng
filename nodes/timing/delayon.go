package timing

import (
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type DelayOn struct {
	*node.Spec
	timer      *time.Timer
	currOutput bool
	lastDelay  time.Duration
}

func NewDelayOn(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOn, category)
	in := node.BuildInput(node.Inp, node.TypeBool, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	delayInput := node.BuildInput(node.Delay, node.TypeFloat, nil, body.Inputs, str.New("interval"))
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(in, delayInput, reset)

	out := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	body.SetSchema(buildDefaultSchema())

	return &DelayOn{body, nil, false, 1 * time.Second}, nil
}

/*
if node input is true
start delay, after the delay set the output to true
*/

func (inst *DelayOn) Process() {
	delayDuration, _ := inst.ReadPinAsTimeSettings(node.Delay)
	if delayDuration != inst.lastDelay {
		inst.setSubtitle(delayDuration)
		inst.lastDelay = delayDuration
	}

	in1, _ := inst.ReadPinAsBool(node.Inp)
	if !in1 { // any time input is false, set output false and cancel any running timers
		inst.WritePinFalse(node.Outp)
		inst.currOutput = false
		if inst.timer != nil {
			inst.timer.Stop()
			inst.timer = nil
		}
		return
	}

	// input is true

	if inst.currOutput { // input is still active, so output is still active, cancel any running timers (for safety)
		inst.WritePinTrue(node.Outp)
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
			inst.WritePinTrue(node.Outp)
			inst.currOutput = true
			inst.timer = nil
		})
	}
	inst.WritePin(node.Outp, inst.currOutput)
}

func (inst *DelayOn) Start() {
	inst.WritePinFalse(node.Outp)
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
