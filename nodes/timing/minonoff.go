package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type MinOnOff struct {
	*node.Spec
	lastInput          bool
	lastReset          bool
	minOnEnabled       bool
	minOffEnabled      bool
	timeOn             int64
	timeOff            int64
	lastMinOnInterval  time.Duration
	lastMinOffInterval time.Duration
}

func NewMinOnOff(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, minOnOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, nil)
	onInterval := node.BuildInput(node.MinOnTime, node.TypeFloat, 0, body.Inputs, str.New("min_on_interval"))
	offInterval := node.BuildInput(node.MinOffTime, node.TypeFloat, 0, body.Inputs, str.New("min_off_interval"))
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in, onInterval, offInterval, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	minOnActive := node.BuildOutput(node.MinOnActive, node.TypeBool, nil, body.Outputs)
	minOffActive := node.BuildOutput(node.MinOffActive, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out, minOnActive, minOffActive)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	node := &MinOnOff{body, false, true, false, false, 0, 0, 1 * time.Second, 1 * time.Second}
	node.SetSchema(node.buildSchema())
	return node, nil
}

func (inst *MinOnOff) Process() {

	minOnIntervalDuration, _ := inst.ReadPinAsTimeSettings(node.MinOnTime)
	minOffIntervalDuration, _ := inst.ReadPinAsTimeSettings(node.MinOffTime)
	if minOnIntervalDuration != inst.lastMinOnInterval || minOffIntervalDuration != inst.lastMinOffInterval {
		inst.setSubtitle(minOnIntervalDuration, minOffIntervalDuration)
		inst.lastMinOnInterval = minOnIntervalDuration
		inst.lastMinOffInterval = minOffIntervalDuration
	}

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
		minOnIntervalSecs := minOnIntervalDuration.Seconds()
		if elapsed >= int64(minOnIntervalSecs) {
			inst.minOnEnabled = false
			inst.WritePinFalse(node.MinOnActive)
			if input {
				inst.WritePinTrue(node.Out)
			}
		}
	} else if inst.minOffEnabled {
		elapsed = time.Now().Unix() - inst.timeOff
		minOffIntervalSecs := minOffIntervalDuration.Seconds()
		if elapsed >= int64(minOffIntervalSecs) {
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

func (inst *MinOnOff) setSubtitle(minOnIntervalDuration, minOffIntervalDuration time.Duration) {
	subtitleText := fmt.Sprintf("min-on %s, min-off %s", minOnIntervalDuration.String(), minOffIntervalDuration.String())
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type MinOnOffSettingsSchema struct {
	MinOnInterval   schemas.Number     `json:"min_on_interval"`
	MinOnTimeUnits  schemas.EnumString `json:"min_on_interval_time_units"`
	MinOffInterval  schemas.Number     `json:"min_off_interval"`
	MinOffTimeUnits schemas.EnumString `json:"min_off_interval_time_units"`
}

type MinOnOffSettings struct {
	MinOnInterval   float64 `json:"min_on_interval"`
	MinOnTimeUnits  string  `json:"min_on_interval_time_units"`
	MinOffInterval  float64 `json:"min_off_interval"`
	MinOffTimeUnits string  `json:"min_off_interval_time_units"`
}

func (inst *MinOnOff) buildSchema() *schemas.Schema {
	props := &MinOnOffSettingsSchema{}
	// time selection
	props.MinOnInterval.Title = "Min ON Interval"
	props.MinOnInterval.Default = 1

	props.MinOffInterval.Title = "Min OFF Interval"
	props.MinOffInterval.Default = 1

	// time selection
	props.MinOnTimeUnits.Title = "Min ON Interval Units"
	props.MinOnTimeUnits.Default = ttime.Sec
	props.MinOnTimeUnits.Options = []string{ttime.Sec, ttime.Min, ttime.Hr}
	props.MinOnTimeUnits.EnumName = []string{ttime.Sec, ttime.Min, ttime.Hr}

	props.MinOffTimeUnits.Title = "Min OFF Interval Units"
	props.MinOffTimeUnits.Default = ttime.Sec
	props.MinOffTimeUnits.Options = []string{ttime.Sec, ttime.Min, ttime.Hr}
	props.MinOffTimeUnits.EnumName = []string{ttime.Sec, ttime.Min, ttime.Hr}

	schema.Set(props)

	uiSchema := array.Map{
		"min_on_interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"min_off_interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"min_on_interval", "min_on_interval_time_units", "min_off_interval", "min_off_interval_time_units"},
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

func (inst *MinOnOff) getSettings(body map[string]interface{}) (*MinOnOffSettings, error) {
	settings := &MinOnOffSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
