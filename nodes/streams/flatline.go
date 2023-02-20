package streams

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

type Flatline struct {
	*node.Spec
	lastVal      *float64
	alertStatus  bool
	lastInterval time.Duration
	lastCOVTime  int64
}

func NewFlatline(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flatLine, category)

	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false) // TODO: this input shouldn't have a manual override value
	delayInput := node.BuildInput(node.Delay, node.TypeFloat, 30, body.Inputs, true)
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

	in, inNull := inst.ReadPinAsFloat(node.In)
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
	Delay          schemas.Number     `json:"delay"`
	DelayTimeUnits schemas.EnumString `json:"delay_time_units"`
}

type FlatlineSettings struct {
	Delay          float64 `json:"delay"`
	DelayTimeUnits string  `json:"delay_time_units"`
}

func (inst *Flatline) buildSchema() *schemas.Schema {
	props := &FlatlineSettingsSchema{}
	// time selection
	props.Delay.Title = "Delay"
	props.Delay.Default = 30

	// time selection
	props.DelayTimeUnits.Title = "Delay Units"
	props.DelayTimeUnits.Default = ttime.Min
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

func (inst *Flatline) getSettings(body map[string]interface{}) (*FlatlineSettings, error) {
	settings := &FlatlineSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
