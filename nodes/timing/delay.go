package timing

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"strings"
	"time"
)

type Delay struct {
	*node.Spec
	activeDelays []*DelayTimer
	lastValue    *float64
}

type DelayTimer struct {
	HasTriggered bool
	Timer        *time.Timer
}

func NewDelay(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, delay, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	delay := node.BuildInput(node.Delay, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(enable, in, delay)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())

	delayArray := make([]*DelayTimer, 0)
	return &Delay{body, delayArray, nil}, nil
}

func (inst *Delay) Process() {

	enable, _ := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.ClearAllDelays()
		return
	}

	input, null := inst.ReadPinAsFloat(node.In)
	inputFloatPtr := float.New(input)
	if null {
		inputFloatPtr = nil
	}

	// if (inputFloatPtr == nil && inst.lastValue != nil) || (inputFloatPtr != nil && inst.lastValue == nil) || *inputFloatPtr != *inst.lastValue {
	if !float.ComparePtrValues(inst.lastValue, inputFloatPtr) {
		settings, _ := getSettings(inst.GetSettings())
		if settings != nil {
			t := strings.Replace(settings.Duration.String(), "ns", "", -1)
			inst.SetSubTitle(fmt.Sprintf("setting: %s %s", t, settings.Time))
		}

		// TODO: Implement delay from wired input

		delayDuration := duration(settings.Duration, settings.Time)

		newDelay := &DelayTimer{false, nil}
		newDelay.Timer = time.AfterFunc(delayDuration, func() {
			delayObj := newDelay
			if inputFloatPtr == nil {
				inst.WritePinNull(node.Out)
			} else {
				inst.WritePinFloat(node.Out, *inputFloatPtr)
			}
			delayObj.HasTriggered = true
		})
		inst.activeDelays = append(inst.activeDelays, newDelay)
	}
	inst.lastValue = inputFloatPtr
	inst.ClearCompletedDelays()
}

func (inst *Delay) ClearCompletedDelays() {
	newDelaysSlice := make([]*DelayTimer, 0)
	for i, _ := range inst.activeDelays {
		if !inst.activeDelays[i].HasTriggered {
			newDelaysSlice = append(newDelaysSlice, inst.activeDelays[i])
		}
	}
	inst.activeDelays = newDelaysSlice
}

func (inst *Delay) ClearAllDelays() {
	for len(inst.activeDelays) > 0 {
		lastIndex := len(inst.activeDelays) - 1
		if inst.activeDelays[lastIndex].Timer != nil {
			inst.activeDelays[lastIndex].Timer.Stop()
			inst.activeDelays[lastIndex] = nil
		}
		inst.activeDelays = inst.activeDelays[:lastIndex]
	}
	inst.activeDelays = make([]*DelayTimer, 0)
}

func (inst *Delay) Stop() {
	inst.ClearAllDelays()
}
