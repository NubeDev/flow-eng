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
	inputs := node.BuildInputs(inSelect)

	outputsCount := int(conversions.GetFloat(body.Settings["outputCount"]))
	dynamicOutputs := node.DynamicOutputs(node.TypeFloat, nil, outputsCount, 2, 20, body.Outputs)
	outputs := node.BuildOutputs(dynamicOutputs...)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetDynamicInputs()
	node := &NumOutputSelect{body}
	node.SetSchema(node.buildSchema())
	return node, nil
}

func (inst *NumOutputSelect) Process() {
	selectInput, selectNull := inst.ReadPinAsFloat(node.Selection)
	if selectNull {
		inst.WritePinNull(node.Outp)
	} else {
		selectInput = math.Floor(selectInput)
		settings, _ := inst.getSettings(inst.GetSettings())
		count := settings.OutputCount
		if selectInput > 0 && int(selectInput) <= count {
			selectedInputName := node.InputName(fmt.Sprintf("in%d", int(selectInput)))
			input, inNull := inst.ReadPinAsFloat(selectedInputName)
			if inNull {
				inst.WritePinNull(node.Outp)
			} else {
				inst.WritePinFloat(node.Outp, input)
			}
		} else {
			inst.WritePinNull(node.Outp)
		}
	}
}

func (inst *NumOutputSelect) SetOutputs(selectIn int, writeValue *float64) {
	settings, _ := inst.getSettings(inst.GetSettings())
	count := settings.OutputCount
	for i := 1; i <= count; i++ {
		selectedOutputName := node.OutputName(fmt.Sprintf("out%d", i))
		if selectIn <= 0 || selectIn > count || selectIn != count || writeValue == nil {
			inst.WritePinNull(selectedOutputName)
		} else {
			inst.WritePinFloat(selectedOutputName, *writeValue)
		}
	}
}

// Custom Node Settings Schema

type NumOutputSelectSettingsSchema struct {
	OutputCount schemas.Integer `json:"outputCount"`
}

type NumOutputSelectSettings struct {
	OutputCount int `json:"outputCount"`
}

func (inst *NumOutputSelect) buildSchema() *schemas.Schema {
	props := &NumOutputSelectSettingsSchema{}

	// inputs count
	props.OutputCount.Title = "Outputs Count"
	props.OutputCount.Default = 2
	props.OutputCount.Minimum = 2
	props.OutputCount.Maximum = 20

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"outputCount"},
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
