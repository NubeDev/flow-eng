package boolean

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

const (
	category = "boolean"
)

const (
	and           = "and"
	or            = "or"
	not           = "not"
	xor           = "xor"
	toggle        = "toggle"
	delayMinOnOff = "min on off"
)

const (
	inputCount = "Inputs Count"
)

type nodeSettings struct {
	InputCount int `json:"inputCount"`
}

func nodeDefault(body *node.Spec, nodeName, category string) (*node.Spec, error) {
	body = node.Defaults(body, nodeName, category)
	settings := &BoolDefaultSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	var count = 2
	if settings != nil {
		count = settings.InputCount
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeBool, false, count, 2, 20, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildBoolDefaultSchema())
	body.SetDynamicInputs()
	return body, nil
}

func Process(body node.Node) {
	equation := body.GetName()
	count := body.InputsLen()
	inputs := boolean.ConvertInterfaceToBoolMultiple(body.ReadMultiple(count))
	body.WritePinBool(node.Out, operation(equation, inputs))
}

func operation(operation string, values []*bool) bool {
	var nonNilValues []bool
	for _, value := range values {
		if value != nil {
			nonNilValues = append(nonNilValues, *value)
		}
	}
	if len(nonNilValues) == 0 {
		return false
	}
	switch operation {
	case and:
		return array.AllTrue(nonNilValues)
	case or:
		return array.OneIsTrue(nonNilValues)
	}
	return false
}

// Custom Node Settings Schema

type BoolDefaultSettingsSchema struct {
	InputCount schemas.Integer `json:"inputCount"`
}

type BoolDefaultSettings struct {
	InputCount int `json:"inputCount"`
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
