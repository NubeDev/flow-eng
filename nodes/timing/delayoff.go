package timing

import (
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
)

type DelayOff struct {
	*node.Spec
	timer   timer.TimedDelay
	wasTrue bool
}

func NewDelayOff(body *node.Spec, timer timer.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delayOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in)
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body.SetSchema(buildSchema())
	return &DelayOff{body, timer, false}, nil
}

func (inst *DelayOff) Process() {
	settings, _ := getSettings(inst.GetSettings())
	in1, null := inst.ReadPinAsBool(node.In)
	if null {
		inst.WritePinNull(node.Out)
	}
	if in1 { // is true
		inst.wasTrue = true
		inst.WritePinTrue(node.Out)
	}

	if !in1 && inst.wasTrue { // was true and now is false
		if inst.timer.WaitFor(duration(settings.Duration, settings.Time)) {
			inst.WritePinFalse(node.Out)
			inst.wasTrue = false
			return
		} else {
			inst.WritePinTrue(node.Out)
			return
		}
	}
	if !in1 {
		inst.WritePinFalse(node.Out)
	}

}
