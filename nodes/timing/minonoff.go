package timing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
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
	currentOutput      bool
	currentMinOn       bool
	currentMinOff      bool
}

func NewMinOnOff(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, minOnOff, Category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false, false)
	onInterval := node.BuildInput(node.MinOnTime, node.TypeFloat, 1, body.Inputs, true, false)
	offInterval := node.BuildInput(node.MinOffTime, node.TypeFloat, 1, body.Inputs, true, false)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in, onInterval, offInterval, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	minOnActive := node.BuildOutput(node.MinOnActive, node.TypeBool, nil, body.Outputs)
	minOffActive := node.BuildOutput(node.MinOffActive, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out, minOnActive, minOffActive)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	n := &MinOnOff{body, false, true, false, false, 0, 0, -1, -1, false, false, false}
	n.SetSchema(n.buildSchema())
	return n, nil
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
		inst.currentMinOn = false
		inst.currentMinOff = false
	}
	inst.lastReset = reset

	if !inst.minOnEnabled && !inst.minOffEnabled {
		inst.WritePinBool(node.Out, input)
		inst.currentOutput = input
		if input && !inst.lastInput {
			inst.timeOn = time.Now().Unix()
			inst.minOnEnabled = true
			inst.WritePinTrue(node.MinOnActive)
			inst.WritePinFalse(node.MinOffActive)
		} else if !input && inst.lastInput {
			inst.timeOff = time.Now().Unix()
			inst.minOffEnabled = true
			inst.WritePinTrue(node.MinOffActive)
			inst.WritePinFalse(node.MinOnActive)
		} else {
			inst.WritePinFalse(node.MinOnActive)
			inst.WritePinFalse(node.MinOffActive)
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
			if input {
				inst.currentOutput = true
			}
		}
	} else if inst.minOffEnabled {
		elapsed = time.Now().Unix() - inst.timeOff
		minOffIntervalSecs := minOffIntervalDuration.Seconds()
		if elapsed >= int64(minOffIntervalSecs) {
			inst.minOffEnabled = false
			if !input {
				inst.currentOutput = false
			}
		}
	}
	inst.lastInput = input
	inst.WritePinBool(node.Out, inst.currentOutput)
	inst.WritePinBool(node.MinOffActive, inst.minOffEnabled)
	inst.WritePinBool(node.MinOnActive, inst.minOnEnabled)

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
	MinOnInterval   schemas.Number     `json:"min-on-time"`
	MinOnTimeUnits  schemas.EnumString `json:"min-on-time_time_units"`
	MinOffInterval  schemas.Number     `json:"min-off-time"`
	MinOffTimeUnits schemas.EnumString `json:"min-off-time_time_units"`
}

type MinOnOffSettings struct {
	MinOnInterval   float64 `json:"min-on-time"`
	MinOnTimeUnits  string  `json:"min-on-time_time_units"`
	MinOffInterval  float64 `json:"min-off-time"`
	MinOffTimeUnits string  `json:"min-off-time_time_units"`
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
		"min-on-time_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"min-off-time_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"min-on-time", "min-on-time_time_units", "min-off-time", "min-off-time_time_units"},
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
