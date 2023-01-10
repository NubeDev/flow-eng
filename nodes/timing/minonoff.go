package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"strconv"
	"time"
)

type MinOnOff struct {
	*node.Spec
	lastInput     bool
	lastReset     bool
	minOnEnabled  bool
	minOffEnabled bool
	timeOn        int64
	timeOff       int64
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
	return &MinOnOff{body, false, true, false, false, 0, 0}, nil
}

func (inst *MinOnOff) Process() {
	// settings, _ := getSettings(inst.GetSettings())

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
		onInterval, _ := inst.ReadPinAsFloat(node.MinOnTime) // TODO: update to settings value and units
		if elapsed >= int64(onInterval) {
			inst.minOnEnabled = false
			inst.WritePinFalse(node.MinOnActive)
			if input {
				inst.WritePinTrue(node.Out)
			}
		}
	} else if inst.minOffEnabled {
		elapsed = time.Now().Unix() - inst.timeOff
		offInterval, _ := inst.ReadPinAsFloat(node.MinOffTime) // TODO: update to settings value and units
		if elapsed >= int64(offInterval) {
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

func (inst *MinOnOff) setSubtitle(settings *OneShotSettings) {
	subtitleText := fmt.Sprintf("min-on %s %s", strconv.FormatFloat(settings.Duration, 'f', -1, 64), settings.TimeUnits)
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type MinOnOffSettingsSchema struct {
	Duration  schemas.Number     `json:"duration"`
	TimeUnits schemas.EnumString `json:"time_units"`
	Retrigger schemas.Boolean    `json:"retrigger"`
}

type MinOnOffSettings struct {
	Duration  float64 `json:"duration"`
	TimeUnits string  `json:"time_units"`
	Retrigger bool    `json:"retrigger"`
}

func (inst *MinOnOff) buildSchema() *schemas.Schema {
	props := &MinOnOffSettingsSchema{}
	// time selection
	props.Duration.Title = "Duration"
	props.Duration.Default = 1

	// time selection
	props.TimeUnits.Title = "Time Units"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.TimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// retrigger selection
	props.Retrigger.Title = "Allow Retrigger"
	props.Retrigger.Default = false

	pprint.PrintJSON(props)
	schema.Set(props)

	fmt.Println(fmt.Sprintf("buildSchema() props: %+v", props))
	pprint.PrintJSON(props)

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
	fmt.Println(fmt.Sprintf("buildSchema() s: %+v", s))
	pprint.PrintJSON(s)
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
