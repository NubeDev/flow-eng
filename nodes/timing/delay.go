package timing

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
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
	enable := node.BuildInput(node.Enable, node.TypeBool, true, body.Inputs, false, false)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	delayInput := node.BuildInput(node.Delay, node.TypeFloat, 1, body.Inputs, true, false)
	inputs := node.BuildInputs(enable, in, delayInput)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	delayArray := make([]*DelayTimer, 0)

	n := &Delay{body, delayArray, nil, 5 * time.Second, nil}
	n.SetSchema(n.buildSchema())
	return n, nil
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

	delayDuration, _ := inst.ReadPinAsTimeSettings(node.Delay)
	if delayDuration != inst.lastDelay {
		inst.setSubtitle(delayDuration)
		inst.lastDelay = delayDuration
	}

	// if (inputFloatPtr == nil && inst.lastValue != nil) || (inputFloatPtr != nil && inst.lastValue == nil) || *inputFloatPtr != *inst.lastValue {
	if !float.ComparePtrValues(inst.lastValue, inputFloatPtr) {

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

// Custom Node Settings Schema

type DelaySettingsSchema struct {
	Delay          schemas.Number     `json:"delay"`
	DelayTimeUnits schemas.EnumString `json:"delay_time_units"`
}

type DelaySettings struct {
	Delay          float64 `json:"delay"`
	DelayTimeUnits string  `json:"delay_time_units"`
}

func (inst *Delay) buildSchema() *schemas.Schema {
	props := &DelaySettingsSchema{}
	// time selection
	props.Delay.Title = "Delay"
	props.Delay.Default = 1

	// time selection
	props.DelayTimeUnits.Title = "Delay Units"
	props.DelayTimeUnits.Default = ttime.Sec
	props.DelayTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.DelayTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"delay", "delay_time_units"},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Node Settings",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

func (inst *Delay) getSettings(body map[string]interface{}) (*DelaySettings, error) {
	settings := &DelaySettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
