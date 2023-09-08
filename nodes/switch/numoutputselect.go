package switches

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type NumOutputSelect struct {
	*node.Spec
}

func NewNumOutputSelect(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, numOutputSelect, Category)
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

	inSelect := node.BuildInput(node.Selection, node.TypeFloat, nil, body.Inputs, false, true)
	input := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, true)
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
	in, inNull := inst.ReadPinAsFloat(node.In)
	settings, _ := inst.getSettings(inst.GetSettings())
	// count := settings.OutputCount
	count := inst.OutputsLen()
	nullNonSelected := settings.NullNonSelected
	for i := 1; i <= count; i++ {
		selectedOutputName := node.OutputName(fmt.Sprintf("out%d", i))
		if nullNonSelected && (selectInput <= 0 || selectInput > float64(count) || selectInput != float64(i) || inNull) {
			inst.WritePinNull(selectedOutputName)
		} else {
			if selectInput == float64(i) {
				inst.WritePinFloat(selectedOutputName, in)
			} else {
				out := inst.GetOutput(selectedOutputName)
				outValue := out.GetValue()
				if outValue == nil {
					inst.WritePinNull(selectedOutputName)
				} else {
					inst.WritePinFloat(selectedOutputName, conversions.GetFloat(outValue))
				}

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
