package timing

import (
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type MinOnOff struct {
	*node.Spec
	lastInput     bool
	lastReset     bool
	minOnEnabled  bool
	minOffEnabled bool
	timeOn        int64
	timeOff       int64
}

func NewMinOnOff(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, minOnOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	onInterval := node.BuildInput(node.MinOnTime, node.TypeFloat, 0, body.Inputs)
	offInterval := node.BuildInput(node.MinOffTime, node.TypeFloat, 0, body.Inputs)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs)
	body.Inputs = node.BuildInputs(in, onInterval, offInterval, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	minOnActive := node.BuildOutput(node.MinOnActive, node.TypeBool, nil, body.Outputs)
	minOffActive := node.BuildOutput(node.MinOffActive, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out, minOnActive, minOffActive)
	return &MinOnOff{body, false, false, false, false, 0, 0}, nil
}

func (inst *MinOnOff) Process() {
	// settings, _ := getSettings(inst.GetSettings())

	input, _ := inst.ReadPinAsBool(node.In)

	reset, _ := inst.ReadPinAsBool(node.Reset)

	if reset && !inst.lastReset {
		inst.minOnEnabled = false
		inst.minOffEnabled = false
		inst.WritePinFalse(node.MinOnActive)
		inst.WritePinFalse(node.MinOffActive)
	}
	inst.lastReset = reset

	if !inst.minOnEnabled && !inst.minOffEnabled {
		inst.WritePinBool(node.Out, input)
		if input && !inst.lastInput {
			inst.timeOn = time.Now().Unix()
			inst.minOnEnabled = true
		} else if !input && inst.lastInput {
			inst.timeOff = time.Now().Unix()
			inst.minOffEnabled = true
		}
		inst.lastInput = input
		return
	}

	var elapsed int64
	if inst.minOnEnabled {
		elapsed = time.Now().Unix() - inst.timeOn
		onInterval, _ := inst.ReadPinAsFloat(node.MinOnTime) // TODO: update to settings value and units
		if elapsed >= int64(onInterval) {
			inst.minOnEnabled = false
			inst.WritePinFalse(node.MinOnActive)
			if input {
				inst.WritePinTrue(node.Out)
			}
		}
	} else if inst.minOffEnabled {
		elapsed = time.Now().Unix() - inst.timeOff
		offInterval, _ := inst.ReadPinAsFloat(node.MinOffTime) // TODO: update to settings value and units
		if elapsed >= int64(offInterval) {
			inst.minOffEnabled = false
			inst.WritePinFalse(node.MinOffActive)
			if !input {
				inst.WritePinFalse(node.Out)
			}
		}
	}
	inst.lastInput = input

}

func (inst *MinOnOff) Stop() {
	inst.minOnEnabled = false
	inst.minOffEnabled = false
}
