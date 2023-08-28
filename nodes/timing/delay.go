package timing

import (
	"encoding/json"
	"fmt"
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
	delayDuration time.Duration
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
	inputs := node.BuildInputs(enable, in)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	delayArray := make([]*DelayTimer, 0)

	n := &Delay{body, delayArray, nil, 5 * time.Second, 0, nil}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Delay) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.init()
	}

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

	if inst.delayDuration != inst.lastDelay {
		inst.lastDelay = inst.delayDuration
	}

	if !float.ComparePtrValues(inst.lastValue, inputFloatPtr) {

		newDelay := &DelayTimer{false, nil}
		newDelay.Timer = time.AfterFunc(inst.delayDuration, func() {
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

func (inst *Delay) init() {
	settings, err := inst.getSettings()
	if err != nil {
		return
	}
	if settings == nil {
		return
	}
	delayTime := time.Duration(settings.Delay)
	delaySetting := settings.DelayTimeUnits
	if delaySetting == ttime.Ms {
		inst.delayDuration = delayTime * time.Millisecond
	}
	if delaySetting == ttime.Sec {
		inst.delayDuration = delayTime * time.Second
	}
	if delaySetting == ttime.Min {
		inst.delayDuration = delayTime * time.Minute
	}
	if delaySetting == ttime.Hr {
		inst.delayDuration = delayTime * time.Hour
	}
	inst.setSubtitle()
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

func (inst *Delay) setSubtitle() {
	settings, err := inst.getSettings()
	if err != nil {
		return
	}
	if settings == nil {
		return
	}
	inst.SetSubTitle(fmt.Sprintf("%s %s", fmt.Sprint(float.RoundTo(settings.Delay, 2)), settings.DelayTimeUnits))
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

func (inst *Delay) getSettings() (*DelaySettings, error) {
	settings := &DelaySettings{}
	marshal, err := json.Marshal(inst.GetSettings())
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
