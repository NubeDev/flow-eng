package count

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/str"
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
}

func NewCount(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, countNode, category)
	countUp := node.BuildInput(node.CountUp, node.TypeBool, nil, body.Inputs, nil)
	countDown := node.BuildInput(node.CountDown, node.TypeBool, nil, body.Inputs, nil)
	stepSize := node.BuildInput(node.StepSize, node.TypeFloat, nil, body.Inputs, str.New("step_size"))
	setValue := node.BuildInput(node.SetValue, node.TypeFloat, nil, body.Inputs, str.New("set_value"))
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, nil)
	body.Inputs = node.BuildInputs(countUp, countDown, stepSize, setValue, reset)

	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)

	node := &Count{body, 0, true, true, true}
	node.SetSchema(node.buildSchema())
	return node, nil
}

func (inst *Count) Process() {
	reset, _ := inst.ReadPinAsBool(node.Reset)
	if reset && !inst.lastReset {
		setVal := inst.ReadPinOrSettingsFloat(node.SetValue)
		inst.count = setVal
	}
	inst.lastReset = reset

	stepSize := inst.ReadPinOrSettingsFloat(node.StepSize)
	if stepSize < 0 {
		stepSize = 1
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

// Custom Node Settings Schema

type CountSettingsSchema struct {
	StepSize schemas.Number `json:"step_size"`
	SetValue schemas.Number `json:"set_value"`
}

type CountSettings struct {
	StepSize float64 `json:"step_size"`
	SetValue float64 `json:"set_value"`
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
		"ui:order": array.Slice{"step_size", "set_value"},
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
