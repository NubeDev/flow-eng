package trigger

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"math"
	"time"

	"github.com/NubeDev/flow-eng/node"
)

type COVNode struct {
	*node.Spec
	lastValue     *float64
	running       bool
	lastInterval  time.Duration
	lastThreshold float64
	currentOutput bool
}

func NewCOVNode(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, COV, Category)
	input := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, true)
	threshold := node.BuildInput(node.Threshold, node.TypeFloat, nil, body.Inputs, true)
	inputs := node.BuildInputs(input, interval, threshold)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp("when ‘input’ changes value, output becomes ‘true’ for ‘interval’ duration, then ‘output’ changes back to ‘false’. For Numeric ‘input’ values, the change of value must be greater than the ‘threshold’ value to trigger the output. Interval value must be equal or larger than 1.")

	n := &COVNode{body, nil, false, 10 * time.Second, -1, false}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *COVNode) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.In)
	covThreshold := inst.ReadPinOrSettingsFloat(node.Threshold)
	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if intervalDuration != inst.lastInterval || covThreshold != inst.lastThreshold {
		inst.setSubtitle(intervalDuration, covThreshold)
		inst.lastInterval = intervalDuration
		inst.lastThreshold = covThreshold
	}
	intervalDuration.Seconds()

	// outputs false if the input is nil or there is no lastValue
	if inputNull || inst.lastValue == nil {
		inst.WritePinBool(node.Out, false)
		inst.currentOutput = false
	} else {
		// call 'writeOutput' when the absolute diff between last two inputs are larger than 'covThreshold' and there are no previous routine running
		diff := math.Abs(input - *inst.lastValue)
		if diff >= covThreshold && !inst.running {
			go writeOutput(inst, intervalDuration)
			inst.running = true
		}
	}
	inst.lastValue = &input
	inst.WritePinBool(node.Out, inst.currentOutput)
}

func writeOutput(inst *COVNode, duration time.Duration) {
	// set the output pin to true for 'duration' period before setting it to false
	// set inst.running to false after routine is finished
	inst.WritePinBool(node.Out, true)
	inst.currentOutput = true
	time.Sleep(duration)
	inst.WritePinBool(node.Out, false)
	inst.currentOutput = false
	inst.running = false
}

func (inst *COVNode) setSubtitle(intervalDuration time.Duration, threshold float64) {
	subtitleText := fmt.Sprintf("theshold:  %f, interval: ", threshold)
	subtitleText += intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type COVNodeSettingsSchema struct {
	Name              schemas.String     `json:"name"`
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	COVThreshold      schemas.Number     `json:"threshold"`
}

type COVNodeSettings struct {
	Name              string  `json:"name"`
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
	COVThreshold      float64 `json:"threshold"`
}

func (inst *COVNode) buildSchema() *schemas.Schema {
	props := &COVNodeSettingsSchema{}

	// name
	props.Name.Title = "Name"
	props.Name.Default = "Change-Of-Value"

	// time selection
	props.Interval.Title = "Interval"
	props.Interval.Default = 1
	props.IntervalTimeUnits.Title = "Interval Units"
	props.IntervalTimeUnits.Default = ttime.Sec
	props.IntervalTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// threshold
	props.COVThreshold.Title = "COV Threshold"
	props.COVThreshold.Default = 0

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"name", "interval", "interval_time_units", "threshold"},
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

func (inst *COVNode) getSettings(body map[string]interface{}) (*COVNodeSettings, error) {
	settings := &COVNodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
