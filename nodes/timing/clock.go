package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type Clock struct {
	*node.Spec
	lastUpdate   int64
	lastInterval time.Duration
	hour         float64
	min          float64
	sec          float64
	timeString   string
	ms           float64
	longString   string
	tzOffset     float64
	unixSecs     float64
}

func NewClock(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, clock, category)
	interval := node.BuildInput(node.Interval, node.TypeFloat, 10, body.Inputs, true, false)
	inputs := node.BuildInputs(interval)

	timeString := node.BuildOutput(node.TimeString, node.TypeString, nil, body.Outputs)
	hour := node.BuildOutput(node.Hour, node.TypeFloat, nil, body.Outputs)
	min := node.BuildOutput(node.Min, node.TypeFloat, nil, body.Outputs)
	sec := node.BuildOutput(node.Sec, node.TypeFloat, nil, body.Outputs)
	ms := node.BuildOutput(node.Ms, node.TypeFloat, nil, body.Outputs)
	longString := node.BuildOutput(node.LongString, node.TypeString, nil, body.Outputs)
	tzOffset := node.BuildOutput(node.TzOffset, node.TypeFloat, nil, body.Outputs)
	unixSecs := node.BuildOutput(node.UnixSecs, node.TypeFloat, nil, body.Outputs)

	outputs := node.BuildOutputs(timeString, hour, min, sec, ms, longString, tzOffset, unixSecs)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &Clock{body, 0, -1, 0, 0, 0, "", 0, "", 0, 0}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Clock) Process() {
	ClockIntervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if ClockIntervalDuration != inst.lastInterval {
		inst.setSubtitle(ClockIntervalDuration)
		inst.lastInterval = ClockIntervalDuration
	}

	currentTime := time.Now()
	if (float64(inst.lastUpdate) + ClockIntervalDuration.Seconds()) <= float64(currentTime.Unix()) {
		inst.lastUpdate = currentTime.Unix()
		inst.hour = float64(currentTime.Hour())
		inst.min = float64(currentTime.Minute())
		inst.sec = float64(currentTime.Second())
		inst.timeString = fmt.Sprintf("%v:%v:%v", inst.hour, inst.min, inst.sec)
		inst.ms = float64(currentTime.Nanosecond() / 1000000)
		inst.longString = currentTime.Format("Monday, January 2, 2006 3:04:05 PM MST")
		_, offsetSeconds := currentTime.Zone()
		inst.tzOffset = float64(offsetSeconds / 3600)
		inst.unixSecs = float64(currentTime.Unix())
	}

	inst.WritePin(node.TimeString, inst.timeString)
	inst.WritePinFloat(node.Hour, inst.hour)
	inst.WritePinFloat(node.Min, inst.min)
	inst.WritePinFloat(node.Sec, inst.sec)
	inst.WritePinFloat(node.Ms, inst.ms)
	inst.WritePin(node.LongString, inst.longString)
	inst.WritePinFloat(node.TzOffset, inst.tzOffset)
	inst.WritePinFloat(node.UnixSecs, inst.unixSecs)
}

func (inst *Clock) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type ClockSettingsSchema struct {
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
}

type ClockSettings struct {
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
}

func (inst *Clock) buildSchema() *schemas.Schema {
	props := &ClockSettingsSchema{}
	// time selection
	props.Interval.Title = "Update Interval"
	props.Interval.Default = 10

	// time selection
	props.IntervalTimeUnits.Title = "Interval Units"
	props.IntervalTimeUnits.Default = ttime.Sec
	props.IntervalTimeUnits.Options = []string{ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Sec, ttime.Min, ttime.Hr}

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "interval_time_units"},
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

func (inst *Clock) getSettings(body map[string]interface{}) (*ClockSettings, error) {
	settings := &ClockSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
