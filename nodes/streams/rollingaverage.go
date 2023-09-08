package streams

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type RollingAverage struct {
	*node.Spec
	lastInterval        time.Duration
	lastSampleTimeMilli int64
	numSamples          int
	lastReset           bool
	sampleArray         []float64
	lastOutput          *float64
}

func NewRollingAverage(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, rollingAverage, Category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, true)
	interval := node.BuildInput(node.Interval, node.TypeFloat, 30, body.Inputs, true, false)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false, true)
	inputs := node.BuildInputs(in, interval, reset)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &RollingAverage{body, 0, 0, 10, true, []float64{}, nil}
	n.SetSchema(n.buildSchema())
	return n, nil

}

func (inst *RollingAverage) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())
	numSamples := settings.NumberOfSamples

	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if intervalDuration != inst.lastInterval || numSamples != inst.numSamples {
		inst.setSubtitle(numSamples, intervalDuration)
		inst.lastInterval = intervalDuration
		inst.numSamples = numSamples

	}

	now := time.Now().UnixMilli()
	inputVal, inNull := inst.ReadPinAsFloat(node.In)

	reset, _ := inst.ReadPinAsBool(node.Reset)
	resetDone := false
	if reset && !inst.lastReset {
		inst.sampleArray = make([]float64, 0)
		inst.lastSampleTimeMilli = now
		resetDone = true
	}
	inst.lastReset = reset

	sampleDelay := intervalDuration.Milliseconds() / int64(inst.numSamples)
	takeSample := resetDone || (now-inst.lastSampleTimeMilli) > sampleDelay

	if !inNull && takeSample {
		// Take a sample and add it to the array
		inst.sampleArray = append(inst.sampleArray, inputVal)
		for len(inst.sampleArray) >= inst.numSamples {
			inst.sampleArray = inst.sampleArray[1:]
		}
		total := 0.0
		for _, value := range inst.sampleArray {
			total += value
		}
		inst.lastOutput = float.New(total / float64(len(inst.sampleArray)))
		inst.lastSampleTimeMilli = now
	}
	if inst.lastOutput == nil {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinFloat(node.Out, *inst.lastOutput)
	}

}

func (inst *RollingAverage) setSubtitle(numSamples int, intervalDuration time.Duration) {
	subtitleText := fmt.Sprintf("%v samples per %s", numSamples, intervalDuration.String())
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type RollingAverageSettingsSchema struct {
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	NumberOfSamples   schemas.Integer    `json:"number_of_samples"`
}

type RollingAverageSettings struct {
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
	NumberOfSamples   int     `json:"number_of_samples"`
}

func (inst *RollingAverage) buildSchema() *schemas.Schema {
	props := &RollingAverageSettingsSchema{}
	// time selection
	props.Interval.Title = "Rolling Average Period"
	props.Interval.Default = 30

	// time selection
	props.IntervalTimeUnits.Title = "Rolling Average Period Units"
	props.IntervalTimeUnits.Default = ttime.Sec
	props.IntervalTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// time selection
	props.NumberOfSamples.Title = "Number Of Samples (per period)"
	props.NumberOfSamples.Default = 30
	props.NumberOfSamples.Minimum = 2

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "interval_time_units", "number_of_samples"},
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

func (inst *RollingAverage) getSettings(body map[string]interface{}) (*RollingAverageSettings, error) {
	settings := &RollingAverageSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
