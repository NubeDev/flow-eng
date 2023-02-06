package numtransform

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"math"
	"time"
)

type RateLimit struct {
	*node.Spec
	lastReset    bool
	lastStep     float64
	lastOutput   float64
	lastInterval time.Duration
}

func NewRateLimit(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, rateLimit, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false)
	input := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false)
	step := node.BuildInput(node.StepSize, node.TypeFloat, nil, body.Inputs, true)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, true)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false)
	inputs := node.BuildInputs(enable, input, step, interval, reset)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &RateLimit{body, true, 0, 0, 1 * time.Second}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *RateLimit) Process() {
	step := inst.ReadPinOrSettingsFloat(node.StepSize)
	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if intervalDuration != inst.lastInterval || step != inst.lastStep {
		inst.setSubtitle(intervalDuration, step)
		inst.lastInterval = intervalDuration
		inst.lastStep = step
	}

	enable, _ := inst.ReadPinAsBool(node.Enable)
	input, inNull := inst.ReadPinAsFloat(node.In)
	reset, _ := inst.ReadPinAsBool(node.Reset)

	output := float64(0)
	if !enable || inNull {
		inst.WritePinFloat(node.Out, 0)
		inst.lastOutput = 0
	} else {
		if reset && !inst.lastReset {
			inst.lastOutput = input
		}
		change := input - inst.lastOutput
		if math.Abs(change) >= step {
			if change < 0 {
				output = inst.lastOutput - step
			} else {
				output = inst.lastOutput + step
			}
		} else {
			output = inst.lastOutput + step
		}
		inst.WritePinFloat(node.Out, output)
		inst.lastOutput = output
	}
}

func (inst *RateLimit) setSubtitle(intervalDuration time.Duration, stepSize float64) {
	subtitleText := fmt.Sprintf("%f every %s", stepSize, intervalDuration.String())
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type RateLimitSettingsSchema struct {
	Interval  schemas.Number     `json:"interval"`
	TimeUnits schemas.EnumString `json:"interval_time_units"`
	StepSize  schemas.Number     `json:"step-size"`
}

type RateLimitSettings struct {
	Interval  float64 `json:"interval"`
	TimeUnits string  `json:"interval_time_units"`
	StepSize  float64 `json:"step-size"`
}

func (inst *RateLimit) buildSchema() *schemas.Schema {
	props := &RateLimitSettingsSchema{}
	// time selection
	props.Interval.Title = "Update Interval"
	props.Interval.Default = 10

	// time selection
	props.TimeUnits.Title = "Update Interval Units"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.TimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// time selection
	props.StepSize.Title = "Maximum Step Size"
	props.StepSize.Default = 1
	props.StepSize.Minimum = 0.00000001

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "interval_time_units", "step-size"},
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

func (inst *RateLimit) getSettings(body map[string]interface{}) (*RateLimitSettings, error) {
	settings := &RateLimitSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
