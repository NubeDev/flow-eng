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

type NumOutputSelect struct {
	*node.Spec
}

func NewNumOutputSelect(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, numOutputSelect, category)
	settings := &NumOutputSelectSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		body.Settings = map[string]interface{}{}
		body.Settings["outputCount"] = 2
	} else if settings.OutputCount < 2 {
		body.Settings["outputCount"] = 2
	}

	inSelect := node.BuildInput(node.Selection, node.TypeFloat, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	input := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)          // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(inSelect, input)

	outputsCount := int(conversions.GetFloat(body.Settings["outputCount"]))
	dynamicOutputs := node.DynamicOutputs(node.TypeFloat, nil, outputsCount, 2, 20, body.Outputs)
	outputs := node.BuildOutputs(dynamicOutputs...)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetDynamicInputs()
	n := &NumOutputSelect{body}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *NumOutputSelect) Process() {
	selectInput, _ := inst.ReadPinAsFloat(node.Selection)
	selectInput = math.Floor(selectInput)
	in, inNull := inst.ReadPinAsFloat(node.Inp)
	settings, _ := inst.getSettings(inst.GetSettings())
	count := settings.OutputCount
	nullNonSelected := settings.NullNonSelected
	for i := 1; i <= count; i++ {
		selectedOutputName := node.OutputName(fmt.Sprintf("out%d", i))
		if nullNonSelected && (selectInput <= 0 || selectInput > float64(count) || selectInput != float64(i) || inNull) {
			inst.WritePinNull(selectedOutputName)
		} else {
			if selectInput == float64(i) {
				inst.WritePinFloat(selectedOutputName, in)
			}
		}
	}
}

// Custom Node Settings Schema

type NumOutputSelectSettingsSchema struct {
	OutputCount     schemas.Integer `json:"outputCount"`
	NullNonSelected schemas.Boolean `json:"null-non-selected"`
}

type NumOutputSelectSettings struct {
	OutputCount     int  `json:"outputCount"`
	NullNonSelected bool `json:"null-non-selected"`
}

func (inst *NumOutputSelect) buildSchema() *schemas.Schema {
	props := &NumOutputSelectSettingsSchema{}

	// inputs count
	props.OutputCount.Title = "Outputs Count"
	props.OutputCount.Default = 2
	props.OutputCount.Minimum = 2
	props.OutputCount.Maximum = 20

	// null non selected
	props.NullNonSelected.Title = "Null Non-Selected Outputs"
	props.NullNonSelected.Default = true

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"outputCount", "null-non-selected"},
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

func (inst *NumOutputSelect) getSettings(body map[string]interface{}) (*NumOutputSelectSettings, error) {
	settings := &NumOutputSelectSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
