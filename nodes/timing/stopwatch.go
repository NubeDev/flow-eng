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

type Stopwatch struct {
	*node.Spec
	lastReset    bool
	lastEnable   bool
	lastUnits    string
	lastTime     int64
	lastOutput   float64
	accumulation float64
}

func NewStopwatch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, stopwatch, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, nil)   // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(enable, reset)

	out := node.BuildOutput(node.Elapsed, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	node := &Stopwatch{body, true, false, "", -1, 0, 0}
	node.SetSchema(node.buildSchema())
	return node, nil
}

func (inst *Stopwatch) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())
	timeUnits := settings.TimeUnits

	if timeUnits != inst.lastUnits {
		inst.setSubtitle(timeUnits)
		inst.lastUnits = timeUnits
	}

	reset, _ := inst.ReadPinAsBool(node.Reset)
	if reset && !inst.lastReset {
		inst.accumulation = 0
		inst.lastOutput = 0
		inst.WritePinFloat(node.Elapsed, 0)
	}
	inst.lastReset = reset

	enable, _ := inst.ReadPinAsBool(node.Enable)
	if enable || (!enable && inst.lastEnable) {
		if !inst.lastEnable {
			inst.lastTime = time.Now().Unix()
		}
		now := time.Now().Unix()
		elapsedSec := float64(now-inst.lastTime) + inst.accumulation
		inst.accumulation = elapsedSec
		var elapsedOut float64
		switch timeUnits {
		case ttime.Ms:
			elapsedOut = elapsedSec * 1000
		case ttime.Sec:
			elapsedOut = elapsedSec
		case ttime.Min:
			elapsedOut = elapsedSec / 60
		case ttime.Hr:
			elapsedOut = (elapsedSec / 60) / 60
		case ttime.Day:
			elapsedOut = ((elapsedSec / 60) / 60) / 24
		}
		inst.lastOutput = elapsedOut
		inst.lastTime = time.Now().Unix()
	}
	inst.lastEnable = enable

	inst.WritePinFloat(node.Elapsed, inst.lastOutput, 4)
}

func (inst *Stopwatch) Start() {
	inst.WritePinFloat(node.Elapsed, 0)
}

func (inst *Stopwatch) setSubtitle(timeUnits string) {
	inst.SetSubTitle(timeUnits)
}

// Custom Node Settings Schema

type StopwatchSettingsSchema struct {
	TimeUnits schemas.EnumString `json:"time_units"`
}

type StopwatchSettings struct {
	TimeUnits string `json:"time_units"`
}

func (inst *Stopwatch) buildSchema() *schemas.Schema {
	props := &StopwatchSettingsSchema{}

	// time selection
	props.TimeUnits.Title = "Interval Units"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr, ttime.Day}
	props.TimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr, ttime.Day}

	schema.Set(props)

	uiSchema := array.Map{
		"time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"time_units"},
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

func (inst *Stopwatch) getSettings(body map[string]interface{}) (*StopwatchSettings, error) {
	settings := &StopwatchSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
