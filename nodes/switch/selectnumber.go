package switches

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
	"math"
)

type SelectNum struct {
	*node.Spec
}

func NewSelectNum(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, selectNum, category)
	settings := &SelectNumSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		body.Settings = map[string]interface{}{}
		body.Settings["inputCount"] = 2
	} else if settings.InputCount < 2 {
		body.Settings["inputCount"] = 2
	}
	inputsCount := int(conversions.GetFloat(body.Settings["inputCount"]))
	inSelect := node.BuildInput(node.Selection, node.TypeFloat, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	dynamicInputs := node.DynamicInputs(node.TypeFloat, nil, inputsCount, 2, 20, body.Inputs)
	inputsArray := []*node.Input{inSelect}
	inputsArray = append(inputsArray, dynamicInputs...)
	inputs := node.BuildInputs(inputsArray...)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetDynamicInputs()
	n := &SelectNum{body}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *SelectNum) Process() {
	selectInput, selectNull := inst.ReadPinAsFloat(node.Selection)
	if selectNull {
		inst.WritePinNull(node.Out)
	} else {
		selectInput = math.Floor(selectInput)
		settings, _ := inst.getSettings(inst.GetSettings())
		count := settings.InputCount
		if selectInput > 0 && int(selectInput) <= count {
			selectedInputName := node.InputName(fmt.Sprintf("in%d", int(selectInput)))
			input, inNull := inst.ReadPinAsFloat(selectedInputName)
			if inNull {
				inst.WritePinNull(node.Out)
			} else {
				inst.WritePinFloat(node.Out, input)
			}
		} else {
			inst.WritePinNull(node.Out)
		}
	}
}

// Custom Node Settings Schema

type SelectNumSettingsSchema struct {
	InputCount schemas.Integer `json:"inputCount"`
}

type SelectNumSettings struct {
	InputCount int `json:"inputCount"`
}

func (inst *SelectNum) buildSchema() *schemas.Schema {
	props := &SelectNumSettingsSchema{}

	// inputs count
	props.InputCount.Title = "Inputs Count"
	props.InputCount.Default = 2
	props.InputCount.Minimum = 2
	props.InputCount.Maximum = 20

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"inputCount"},
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

func (inst *SelectNum) getSettings(body map[string]interface{}) (*SelectNumSettings, error) {
	settings := &SelectNumSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
