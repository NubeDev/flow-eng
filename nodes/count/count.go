package count

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type Count struct {
	*node.Spec
	count         float64
	lastReset     bool
	lastCountUp   bool
	lastCountDown bool
	lastStepSize  float64
	lastSetVal    float64
}

func NewCount(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, countNode, category)
	countUp := node.BuildInput(node.CountUp, node.TypeBool, nil, body.Inputs, false, true)
	countDown := node.BuildInput(node.CountDown, node.TypeBool, nil, body.Inputs, false, true)
	stepSize := node.BuildInput(node.StepSize, node.TypeFloat, 1, body.Inputs, true, false)
	setValue := node.BuildInput(node.SetValue, node.TypeFloat, 0, body.Inputs, true, false)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false, true)
	body.Inputs = node.BuildInputs(countUp, countDown, stepSize, setValue, reset)

	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)

	n := &Count{body, 0, true, true, true, 0, 0}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Count) Process() {
	reset, _ := inst.ReadPinAsBool(node.Reset)
	setVal := inst.ReadPinOrSettingsFloat(node.SetValue)
	if reset && !inst.lastReset {
		inst.count = setVal
	}
	inst.lastReset = reset

	stepSize := inst.ReadPinOrSettingsFloat(node.StepSize)
	if stepSize < 0 {
		stepSize = 1
	}

	if stepSize != inst.lastStepSize || setVal != inst.lastSetVal {
		inst.setSubtitle(stepSize, setVal)
		inst.lastStepSize = stepSize
		inst.lastSetVal = setVal
	}

	countUp, _ := inst.ReadPinAsBool(node.CountUp)
	if countUp && !inst.lastCountUp {
		inst.count += stepSize
	}
	inst.lastCountUp = countUp

	countDown, _ := inst.ReadPinAsBool(node.CountDown)
	if countDown && !inst.lastCountDown {
		inst.count -= stepSize
	}
	inst.lastCountDown = countDown

	inst.WritePinFloat(node.CountOut, inst.count)
}

func (inst *Count) setSubtitle(step, set float64) {
	subtitleText := fmt.Sprintf("step-size %v   set-value: %v", step, set)
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type CountSettingsSchema struct {
	StepSize schemas.Number `json:"step-size"`
	SetValue schemas.Number `json:"set-value"`
}

type CountSettings struct {
	StepSize float64 `json:"step-size"`
	SetValue float64 `json:"set-value"`
}

func (inst *Count) buildSchema() *schemas.Schema {
	props := &CountSettingsSchema{}

	// Step Size
	props.StepSize.Title = "Step Size"
	props.StepSize.Default = 1

	// Set Value
	props.SetValue.Title = "Set Value"
	props.SetValue.Default = 0

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"step-size", "set-value"},
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

func (inst *Count) getSettings(body map[string]interface{}) (*CountSettings, error) {
	settings := &CountSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
