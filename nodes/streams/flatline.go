package streams

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type Flatline struct {
	*node.Spec
	lastVal      *float64
	alertStatus  bool
	lastInterval time.Duration
	lastCOVTime  int64
}

func NewFlatline(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flatLine, category)

	in := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	delayInput := node.BuildInput(node.Delay, node.TypeFloat, nil, body.Inputs, str.New("interval"))
	inputs := node.BuildInputs(in, delayInput)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &Flatline{body, nil, false, -1, -1}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Flatline) Process() {
	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Delay)
	if intervalDuration != inst.lastInterval {
		inst.setSubtitle(intervalDuration)
		inst.lastInterval = intervalDuration
	}

	newCOV := false

	in, inNull := inst.ReadPinAsFloat(node.Inp)
	if inNull && inst.lastVal != nil || inNull && inst.lastCOVTime == -1 {
		newCOV = true
		inst.lastVal = nil
	} else if !inNull && inst.lastVal == nil {
		newCOV = true
		inst.lastVal = float.New(in)
	} else if !inNull && (in != *inst.lastVal) {
		newCOV = true
		inst.lastVal = float.New(in)
	}

	if newCOV {
		inst.lastCOVTime = time.Now().Unix()
		inst.alertStatus = false
	} else {
		now := time.Now().Unix()
		if float64(now-inst.lastCOVTime) >= intervalDuration.Seconds() {
			inst.alertStatus = true
		}
	}
	inst.WritePinBool(node.Out, inst.alertStatus)
}

func (inst *Flatline) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type FlatlineSettingsSchema struct {
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
}

type FlatlineSettings struct {
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
}

func (inst *Flatline) buildSchema() *schemas.Schema {
	props := &FlatlineSettingsSchema{}
	// time selection
	props.Interval.Title = "Interval"
	props.Interval.Default = 30

	// time selection
	props.IntervalTimeUnits.Title = "Interval Units"
	props.IntervalTimeUnits.Default = ttime.Min
	props.IntervalTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

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

func (inst *Flatline) getSettings(body map[string]interface{}) (*FlatlineSettings, error) {
	settings := &FlatlineSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
