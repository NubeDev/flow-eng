package timing

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type OneShot struct {
	*node.Spec
	timer        *time.Timer
	outputActive bool
	lastTrigger  bool
	lastReset    bool
	lastInterval time.Duration
}

func NewOneShot(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, oneShot, category)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, false) // TODO: this input shouldn't have a manual override value
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false)          // TODO: this input shouldn't have a manual override value
	interval := node.BuildInput(node.Interval, node.TypeFloat, 1, body.Inputs, true)
	inputs := node.BuildInputs(trigger, reset, interval)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &OneShot{body, nil, false, true, true, 0}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *OneShot) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())
	retrigger := settings.Retrigger

	oneShotIntervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if oneShotIntervalDuration != inst.lastInterval {
		inst.setSubtitle(oneShotIntervalDuration)
		inst.lastInterval = oneShotIntervalDuration
	}

	trigger, _ := inst.ReadPinAsBool(node.TriggerInput)
	if trigger && !inst.lastTrigger {
		if retrigger || !inst.outputActive {
			inst.StartOneShot(oneShotIntervalDuration)
		}
	}
	inst.lastTrigger = trigger

	reset, _ := inst.ReadPinAsBool(node.Reset)
	if reset && !inst.lastReset {
		if inst.outputActive {
			inst.StopOneShotTimer(true)
		}
	}
	inst.lastReset = reset
	inst.WritePinBool(node.Out, inst.outputActive)
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

func (inst *OneShot) Start() {
	inst.WritePinFalse(node.Out)
	inst.outputActive = false
}

func (inst *OneShot) Stop() {
	inst.StopOneShotTimer(true)
}

func (inst *OneShot) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type OneShotSettingsSchema struct {
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	Retrigger         schemas.Boolean    `json:"retrigger"`
}

type OneShotSettings struct {
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
	Retrigger         bool    `json:"retrigger"`
}

func (inst *OneShot) buildSchema() *schemas.Schema {
	props := &OneShotSettingsSchema{}
	// time selection
	props.Interval.Title = "Interval"
	props.Interval.Default = 1

	// time selection
	props.IntervalTimeUnits.Title = "Interval Units"
	props.IntervalTimeUnits.Default = ttime.Sec
	props.IntervalTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// retrigger selection
	props.Retrigger.Title = "Retrigger"
	props.Retrigger.Default = false
	props.Retrigger.EnumNames = []string{"Retrigger While Output Is Active", "Only Retrigger While Output Is Inactive"}

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"retrigger": array.Map{
			"ui:widget": "select",
		},
		"ui:order": array.Slice{"interval", "interval_time_units", "retrigger"},
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

func (inst *OneShot) getSettings(body map[string]interface{}) (*OneShotSettings, error) {
	settings := &OneShotSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
