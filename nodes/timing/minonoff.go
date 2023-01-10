package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"strconv"
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
	lastMinOnInterval  float64
	lastMinOffInterval float64
}

func NewMinOnOff(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, minOnOff, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	onInterval := node.BuildInput(node.MinOnTime, node.TypeFloat, 0, body.Inputs)
	offInterval := node.BuildInput(node.MinOffTime, node.TypeFloat, 0, body.Inputs)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(in, onInterval, offInterval, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	minOnActive := node.BuildOutput(node.MinOnActive, node.TypeBool, nil, body.Outputs)
	minOffActive := node.BuildOutput(node.MinOffActive, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out, minOnActive, minOffActive)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &MinOnOff{body, false, true, false, false, 0, 0, 0, 0}, nil
}

func (inst *MinOnOff) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())

	// Select intervals from input links or node settings
	minOnIntervalLink, onIntNull := inst.ReadPinAsFloat(node.MinOnTime)
	var minOnIntervalDuration time.Duration
	var minOnInterval float64
	if onIntNull { // no interval input, so get value from settings
		minOnInterval = settings.MinOnInterval
		minOnIntervalDuration = ttime.Duration(settings.MinOnInterval, settings.MinOnTimeUnits)
	} else {
		minOnInterval = minOnIntervalLink
		minOnIntervalDuration = ttime.Duration(minOnIntervalLink, settings.MinOnTimeUnits)
	}
	minOffIntervalLink, offIntNull := inst.ReadPinAsFloat(node.MinOffTime)
	var minOffIntervalDuration time.Duration
	var minOffInterval float64
	if offIntNull { // no interval input, so get value from settings
		minOffInterval = settings.MinOffInterval
		minOffIntervalDuration = ttime.Duration(settings.MinOffInterval, settings.MinOffTimeUnits)
	} else {
		minOffInterval = minOffIntervalLink
		minOffIntervalDuration = ttime.Duration(minOffIntervalLink, settings.MinOffTimeUnits)
	}
	if minOnInterval != inst.lastMinOnInterval || minOffInterval != inst.lastMinOffInterval {
		inst.setSubtitle(minOnInterval, settings.MinOnTimeUnits, minOffInterval, settings.MinOffTimeUnits)
		inst.lastMinOnInterval = minOnInterval
		inst.lastMinOffInterval = minOffInterval
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

func (inst *MinOnOff) setSubtitle(minOnIntervalAmount float64, minOnTimeUnits string, minOffIntervalAmount float64, minOffTimeUnits string) {
	subtitleText := fmt.Sprintf("min-on %s %s, min-off %s %s", strconv.FormatFloat(minOnIntervalAmount, 'f', -1, 64), minOnTimeUnits, strconv.FormatFloat(minOffIntervalAmount, 'f', -1, 64), minOffTimeUnits)
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type MinOnOffSettingsSchema struct {
	MinOnInterval   schemas.Number     `json:"min_on_interval"`
	MinOnTimeUnits  schemas.EnumString `json:"min_on_time_units"`
	MinOffInterval  schemas.Number     `json:"min_off_interval"`
	MinOffTimeUnits schemas.EnumString `json:"min_off_time_units"`
}

type MinOnOffSettings struct {
	MinOnInterval   float64 `json:"min_on_interval"`
	MinOnTimeUnits  string  `json:"min_on_time_units"`
	MinOffInterval  float64 `json:"min_off_interval"`
	MinOffTimeUnits string  `json:"min_off_time_units"`
}

func (inst *MinOnOff) buildSchema() *schemas.Schema {
	props := &MinOnOffSettingsSchema{}
	// time selection
	props.MinOnInterval.Title = "Min ON Interval"
	props.MinOnInterval.Default = 1

	props.MinOffInterval.Title = "Min OFF Interval"
	props.MinOffInterval.Default = 1

	// time selection
	props.MinOnTimeUnits.Title = "Min ON Time Units"
	props.MinOnTimeUnits.Default = ttime.Sec
	props.MinOnTimeUnits.Options = []string{ttime.Sec, ttime.Min, ttime.Hr}
	props.MinOnTimeUnits.EnumName = []string{ttime.Sec, ttime.Min, ttime.Hr}

	props.MinOffTimeUnits.Title = "Min OFF Time Units"
	props.MinOffTimeUnits.Default = ttime.Sec
	props.MinOffTimeUnits.Options = []string{ttime.Sec, ttime.Min, ttime.Hr}
	props.MinOffTimeUnits.EnumName = []string{ttime.Sec, ttime.Min, ttime.Hr}

	schema.Set(props)

	uiSchema := array.Map{
		"time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
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
