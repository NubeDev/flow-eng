package timing

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"strconv"
	"time"
)

type DelayOn struct {
	*node.Spec
	timer      *time.Timer
	currOutput bool
}

func NewDelayOn(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOn, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(in)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildDefaultSchema())
	return &DelayOn{body, nil, false}, nil
}

/*
if node input is true
start delay, after the delay set the output to true
*/

func (inst *DelayOn) Process() {
	settings, _ := getDefaultSettings(inst.GetSettings())
	if settings != nil {
		subtitleText := fmt.Sprintf("%s %s", strconv.FormatFloat(settings.Duration, 'f', -1, 64), settings.TimeUnits)
		inst.SetSubTitle(subtitleText)
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
		onDelayDuration := ttime.Duration(settings.Duration, settings.TimeUnits)
		inst.timer = time.AfterFunc(onDelayDuration, func() {
			inst.WritePinTrue(node.Out)
			inst.currOutput = true
			inst.timer = nil
		})
	}
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
