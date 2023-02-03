package math

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

const (
	category = "math"
	divide   = "divide"
	add      = "add"
	sub      = "subtract"
	multiply = "multiply"
	modulo   = "modulo"
	absolute = "absolute-value"
	floor    = "floor"
	ceiling  = "ceiling"
)

const (
	mathHelp string = "perform a math operation: "
)

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	settings := &MathDefaultSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	var count = 2
	if settings != nil {
		count = settings.InputCount
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, count, 2, 20, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildMathDefaultSchema())
	body.SetDynamicInputs()
	return body, nil
}

func process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := conversions.ConvertInterfaceToFloatMultiple(body.ReadMultiple(count))
	output := operation(equation, inputs)
	if output == nil {
		body.WritePinNull(node.Outp)
	} else {
		body.WritePinFloat(node.Outp, float.NonNil(output))
	}
}

// Custom Node Settings Schema

type MathDefaultSettingsSchema struct {
	InputCount schemas.Integer `json:"inputCount"`
}

type MathDefaultSettings struct {
	InputCount int `json:"inputCount"`
}

func buildMathDefaultSchema() *schemas.Schema {
	props := &MathDefaultSettingsSchema{}
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
