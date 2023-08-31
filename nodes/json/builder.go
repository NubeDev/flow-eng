package nodejson

import (
	"encoding/json"
	"fmt"

	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type JSONNumberBuilder struct {
	*node.Spec
	outputType string
}

const maxInputs = 10

func NewJSONNumberBuilder(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, jsonNumberBuilder, Category)
	settings := &BoolDefaultSettings{}
	mapstructure.Decode(body.Settings, &settings)
	var count = 2
	var outputType string
	if settings != nil {
		count = settings.InputCount
		outputType = settings.OutputType

	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeNumber, nil, count, 2, maxInputs, body.Inputs, false)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildBoolDefaultSchema())
	body.SetDynamicInputs()
	return &JSONNumberBuilder{body, outputType}, nil
}

func (inst *JSONNumberBuilder) Process() {

	inputs := inst.ReadMultipleInputs(maxInputs)
	var values [10]*float64

	for _, input := range inputs {
		switch input.Name {
		case node.In1:
			values[0] = conversions.GetFloatPointer(input.Value)
		case node.In2:
			values[1] = conversions.GetFloatPointer(input.Value)
		case node.In3:
			values[2] = conversions.GetFloatPointer(input.Value)
		case node.In4:
			values[3] = conversions.GetFloatPointer(input.Value)
		case node.In5:
			values[4] = conversions.GetFloatPointer(input.Value)
		case node.In6:
			values[5] = conversions.GetFloatPointer(input.Value)
		case node.In7:
			values[6] = conversions.GetFloatPointer(input.Value)
		case node.In8:
			values[7] = conversions.GetFloatPointer(input.Value)
		case node.In9:
			values[8] = conversions.GetFloatPointer(input.Value)
		case node.In10:
			values[9] = conversions.GetFloatPointer(input.Value)
		}
	}

	if inst.outputType == outputTypeMap {
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
			return
		} else {
			inst.WritePin(node.Out, payload)
			return
		}
	}

	payload, err := json2Str(values)
	if err != nil {
		inst.WritePin(node.Out, err.Error())
		return
	} else {
		inst.WritePin(node.Out, payload)
		return
	}

}

type JSONBuilder struct {
	*node.Spec
	outputType string
}

// NewJSONBuilder for strings
func NewJSONBuilder(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, jsonBuilder, Category)
	settings := &BoolDefaultSettings{}
	mapstructure.Decode(body.Settings, &settings)
	var count = 2
	var outputType string
	if settings != nil {
		count = settings.InputCount
		outputType = settings.OutputType

	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeString, nil, count, 2, 10, body.Inputs, false)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildBoolDefaultSchema())
	body.SetDynamicInputs()
	return &JSONBuilder{body, outputType}, nil
}

func (inst *JSONBuilder) Process() {
	inputs := inst.ReadMultipleInputs(maxInputs)
	var values [10]*string

	for _, input := range inputs {
		switch input.Name {
		case node.In1:
			values[0] = conversions.GetStringPointer(input.Value)
		case node.In2:
			values[1] = conversions.GetStringPointer(input.Value)
		case node.In3:
			values[2] = conversions.GetStringPointer(input.Value)
		case node.In4:
			values[3] = conversions.GetStringPointer(input.Value)
		case node.In5:
			values[4] = conversions.GetStringPointer(input.Value)
		case node.In6:
			values[5] = conversions.GetStringPointer(input.Value)
		case node.In7:
			values[6] = conversions.GetStringPointer(input.Value)
		case node.In8:
			values[7] = conversions.GetStringPointer(input.Value)
		case node.In9:
			values[8] = conversions.GetStringPointer(input.Value)
		case node.In10:
			values[9] = conversions.GetStringPointer(input.Value)
		}
	}

	if inst.outputType == outputTypeMap {
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
			return
		} else {
			inst.WritePin(node.Out, payload)
			return
		}
	}

	payload, err := json2Str(values)
	if err != nil {
		inst.WritePin(node.Out, err.Error())
		return
	} else {
		inst.WritePin(node.Out, payload)
		return
	}

}

type BoolDefaultSettings struct {
	OutputType string `json:"outputType"`
	InputCount int    `json:"inputCount"`
}

type BoolDefaultSettingsSchema struct {
	OutputType schemas.EnumString `json:"outputType"`
	InputCount schemas.Integer    `json:"inputCount"`
}

const (
	outputTypeMap   = "object"
	outputTypeArray = "array"
)

func buildBoolDefaultSchema() *schemas.Schema {
	props := &BoolDefaultSettingsSchema{}

	props.OutputType.Title = "output type"
	props.OutputType.Default = outputTypeMap
	props.OutputType.Options = []string{outputTypeMap, outputTypeArray}
	props.OutputType.EnumName = []string{outputTypeMap, outputTypeArray}

	props.InputCount.Title = "Input Count"
	props.InputCount.Default = 2
	props.InputCount.Minimum = 2
	props.InputCount.Maximum = 10

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"outputType", "inputCount"},
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
