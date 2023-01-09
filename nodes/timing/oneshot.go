package timing

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"strings"
	"time"
)

type OneShot struct {
	*node.Spec
	timer        *time.Timer
	outputActive bool
	lastIn       bool
	lastReset    bool
}

func NewOneShot(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, oneShot, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)       // TODO: this input shouldn't have a manual override value
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(in, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &OneShot{body, nil, false, true, true}, nil
}

func (inst *OneShot) Process() {
	retrigger := false

	settings, _ := getSettings(inst.GetSettings())
	if settings != nil {
		t := strings.Replace(settings.Duration.String(), "ns", "", -1)
		inst.SetSubTitle(fmt.Sprintf("setting: %s %s", t, settings.Time))
	}

	in, _ := inst.ReadPinAsBool(node.In)
	if in && !inst.lastIn {
		if retrigger || !inst.outputActive {
			oneShotDuration := duration(settings.Duration, settings.Time)
			inst.StartOneShot(oneShotDuration)
		}
	}
	inst.lastIn = in

}

func (inst *OneShot) StartOneShot(duration time.Duration) {
	if inst.timer != nil {
		inst.StopOneShotTimer(false)
	}
	inst.timer = time.AfterFunc(duration, func() {
		inst.WritePinFalse(node.Out)
		inst.outputActive = false
		inst.timer = nil
	})
	inst.WritePinTrue(node.Out)
	inst.outputActive = true
}

func (inst *OneShot) StopOneShotTimer(reset bool) {
	if inst.timer != nil {
		inst.timer.Stop()
	}
	if reset {
		inst.WritePinFalse(node.Out)
		inst.outputActive = false
	}
}
