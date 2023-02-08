package timing

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type Delay struct {
	*node.Spec
	activeDelays  []*DelayTimer
	lastValue     *float64
	lastDelay     time.Duration
	currentOutput *float64
}

type DelayTimer struct {
	HasTriggered bool
	Timer        *time.Timer
}

func NewDelay(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, delay, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false)
	delayInput := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, true)
	inputs := node.BuildInputs(enable, in, delayInput)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildDefaultSchema())
	delayArray := make([]*DelayTimer, 0)
	return &Delay{body, delayArray, nil, 1 * time.Second, nil}, nil
}

func (inst *Delay) Process() {
	enable, _ := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.ClearAllDelays()
		inst.WritePinNull(node.Out)
		inst.currentOutput = nil
		return
	}

	input, null := inst.ReadPinAsFloat(node.In)
	inputFloatPtr := float.New(input)
	if null {
		inputFloatPtr = nil
	}

	// if (inputFloatPtr == nil && inst.lastValue != nil) || (inputFloatPtr != nil && inst.lastValue == nil) || *inputFloatPtr != *inst.lastValue {
	if !float.ComparePtrValues(inst.lastValue, inputFloatPtr) {
		delayDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
		if delayDuration != inst.lastDelay {
			inst.setSubtitle(delayDuration)
			inst.lastDelay = delayDuration
		}

		newDelay := &DelayTimer{false, nil}
		newDelay.Timer = time.AfterFunc(delayDuration, func() {
			delayObj := newDelay
			if inputFloatPtr == nil {
				inst.WritePinNull(node.Out)
				inst.currentOutput = nil
			} else {
				inst.WritePinFloat(node.Out, *inputFloatPtr)
				inst.currentOutput = float.New(*inputFloatPtr)
			}
			delayObj.HasTriggered = true
		})
		inst.activeDelays = append(inst.activeDelays, newDelay)
	}
	inst.lastValue = inputFloatPtr
	if inst.currentOutput == nil {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinFloat(node.Out, *inst.currentOutput)
	}
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

func (inst *Delay) Start() {
	inst.WritePinNull(node.Out)
}

func (inst *Delay) Stop() {
	inst.ClearAllDelays()
}

func (inst *Delay) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}
