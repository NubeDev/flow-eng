package switches

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"math"
)

type SelectNum struct {
	*node.Spec
}

func NewSelectNum(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, selectNum, category)
	/*
		buildCount, setting, value, err := node.NewSetting(body, &node.SettingOptions{Type: node.Number, Title: node.InputCount, Min: 2, Max: 20})
		if err != nil {return nil, err
		}
		settings, err := node.BuildSettings(setting)
		if err != nil {return nil, err
		}
		count, ok := value.(int)
		if !ok {count = 2
		}

		node.DynamicInputs(node.TypeFloat, nil, count, 2, 20, body.Inputs)...

		inSelect := node.BuildInput(node.Selection, node.TypeFloat, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
		// inputs := node.BuildInputs(inSwitch, inTrue, inFalse)
		inputs := node.BuildInputs(inSelect, node.DynamicInputs(node.TypeFloat, nil, count, 2, 20, body.Inputs)...)

		out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
		outputs := node.BuildOutputs(out)
		body = node.BuildNode(body, inputs, outputs, body.Settings)
	*/
	return &Switch{body}, nil
}

func (inst *SelectNum) Process() {
	selectInput, selectNull := inst.ReadPinAsFloat(node.Selection)
	if selectNull {
		inst.WritePinNull(node.Outp)
	} else {
		selectInput = math.Floor(selectInput)

	}
}

// Custom Node Settings Schema

type SelectNumSettingsSchema struct {
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	Retrigger         schemas.Boolean    `json:"retrigger"`
}

type SelectNumSettings struct {
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
	Retrigger         bool    `json:"retrigger"`
}

func (inst *SelectNum) buildSchema() *schemas.Schema {
	props := &SelectNumSettingsSchema{}

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"retrigger": array.Map{
			"ui:widget": "select",
		},
		"ui:order": array.Slice{"interval", "interval_time_units", "retrigger"},
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
