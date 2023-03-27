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

type Date struct {
	*node.Spec
	lastUpdate   int64
	lastInterval time.Duration
	dateString   string
	dayString    string
	dow          float64
	dateNum      float64
	month        float64
	year         float64
}

func NewDate(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, date, category)
	interval := node.BuildInput(node.Interval, node.TypeFloat, 60, body.Inputs, true, false)
	inputs := node.BuildInputs(interval)

	dateString := node.BuildOutput(node.DateString, node.TypeString, nil, body.Outputs)
	dayString := node.BuildOutput(node.DayString, node.TypeString, nil, body.Outputs)
	dow := node.BuildOutput(node.DayOfWeek, node.TypeFloat, nil, body.Outputs)
	dateNum := node.BuildOutput(node.Date, node.TypeFloat, nil, body.Outputs)
	month := node.BuildOutput(node.Month, node.TypeFloat, nil, body.Outputs)
	year := node.BuildOutput(node.Year, node.TypeFloat, nil, body.Outputs)

	outputs := node.BuildOutputs(dateString, dayString, dow, dateNum, month, year)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &Date{body, 0, -1, "", "", 0, 0, 0, 0}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Date) Process() {
	DateIntervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if DateIntervalDuration != inst.lastInterval {
		inst.setSubtitle(DateIntervalDuration)
		inst.lastInterval = DateIntervalDuration
	}

	currentTime := time.Now()
	if (float64(inst.lastUpdate) + DateIntervalDuration.Seconds()) <= float64(currentTime.Unix()) {
		inst.lastUpdate = currentTime.Unix()
		inst.dateString = currentTime.Format("02/01/2006")
		inst.dayString = currentTime.Format("Monday")
		inst.dow = float64(currentTime.Weekday())
		inst.dateNum = float64(currentTime.Day())
		inst.month = float64(currentTime.Month())
		inst.year = float64(currentTime.Year())
	}

	inst.WritePin(node.DateString, inst.dateString)
	inst.WritePin(node.DayString, inst.dayString)
	inst.WritePinFloat(node.DayOfWeek, inst.dow)
	inst.WritePinFloat(node.Date, inst.dateNum)
	inst.WritePinFloat(node.Month, inst.month)
	inst.WritePinFloat(node.Year, inst.year)
}

func (inst *Date) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type DateSettingsSchema struct {
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
}

type DateSettings struct {
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
}

func (inst *Date) buildSchema() *schemas.Schema {
	props := &DateSettingsSchema{}
	// time selection
	props.Interval.Title = "Update Interval"
	props.Interval.Default = 60

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

func (inst *Date) getSettings(body map[string]interface{}) (*DateSettings, error) {
	settings := &DateSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
