package nodejson

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type JSONBuilder struct {
	*node.Spec
}

type BoolDefaultSettings struct {
	InputCount int `json:"inputCount"`
}

func NewJSONBuilder(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, jsonBuilder, category)
	settings := &BoolDefaultSettings{}
	mapstructure.Decode(body.Settings, &settings)
	var count = 2
	if settings != nil {
		count = settings.InputCount
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeString, nil, count, 2, 10, body.Inputs, false)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildBoolDefaultSchema())
	body.SetDynamicInputs()
	return &JSONBuilder{body}, nil
}

func (inst *JSONBuilder) Process() {
	count := inst.InputsLen()
	inputs := str.ConvertInterfaceToStringMultiple(inst.ReadMultiple(count))
	var values [10]*string
	for i, v := range inputs {
		values[i] = v
	}
	o := out{
		In1:  values[0],
		In2:  values[1],
		In3:  values[2],
		In4:  values[3],
		In5:  values[4],
		In6:  values[5],
		In7:  values[6],
		In8:  values[7],
		In9:  values[8],
		In10: values[9],
	}

	payload, err := json2Str(o)
	if err != nil {
		inst.WritePin(node.Out, err.Error())
	} else {
		inst.WritePin(node.Out, payload)
	}

}

type BoolDefaultSettingsSchema struct {
	InputCount schemas.Integer `json:"inputCount"`
}

func buildBoolDefaultSchema() *schemas.Schema {
	props := &BoolDefaultSettingsSchema{}
	props.InputCount.Title = "Input Count"
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

type out struct {
	In1  interface{} `json:"in1"`
	In2  interface{} `json:"in2"`
	In3  interface{} `json:"in3"`
	In4  interface{} `json:"in4"`
	In5  interface{} `json:"in5"`
	In6  interface{} `json:"in6"`
	In7  interface{} `json:"in7"`
	In8  interface{} `json:"in8"`
	In9  interface{} `json:"in9"`
	In10 interface{} `json:"in10"`
}

func json2Str(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	return string(b), nil
}
