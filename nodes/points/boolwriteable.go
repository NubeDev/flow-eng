package point

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type BooleanWriteable struct {
	*node.Spec
}

func NewBooleanWriteable(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, booleanWriteable, Category)
	settings := &BooleanWriteableSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	var count = 2
	if settings != nil {
		count = settings.InputCount
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeBool, nil, count, 2, 20, body.Inputs, false)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildBooleanWriteableSchema())
	body.SetDynamicInputs()
	return &BooleanWriteable{body}, nil
}

func (inst *BooleanWriteable) Process() {
	count := inst.InputsLen()
	inputs := boolean.ConvertInterfaceToBoolMultiple(inst.ReadMultiple(count))
	for _, val := range inputs {
		if val != nil {
			inst.WritePinBool(node.Out, *val)
			return
		}
	}
	inst.WritePinNull(node.Out)
}

// Custom Node Settings Schema

type BooleanWriteableSettingsSchema struct {
	InputCount schemas.Integer `json:"inputCount"`
}

type BooleanWriteableSettings struct {
	InputCount int `json:"inputCount"`
}

func buildBooleanWriteableSchema() *schemas.Schema {
	props := &BooleanWriteableSettingsSchema{}
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
