package timing

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"strconv"
	"time"
)

type DelayOff struct {
	*node.Spec
	timer      *time.Timer
	currOutput bool
}

func NewDelayOff(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	body.SetSchema(buildDefaultSchema())
	return &DelayOff{body, nil, false}, nil
}

func (inst *DelayOff) Process() {
	settings, _ := getDefaultSettings(inst.GetSettings())
	if settings != nil {
		subtitleText := fmt.Sprintf("%s %s", strconv.FormatFloat(settings.Duration, 'f', -1, 64), settings.TimeUnits)
		inst.SetSubTitle(subtitleText)
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
		onDelayDuration := ttime.Duration(settings.Duration, settings.TimeUnits)
		inst.timer = time.AfterFunc(onDelayDuration, func() {
			inst.WritePinFalse(node.Out)
			inst.currOutput = false
			inst.timer = nil
		})
	}
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
