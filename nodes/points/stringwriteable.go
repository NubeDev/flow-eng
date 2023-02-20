package point

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type StringWriteable struct {
	*node.Spec
}

func NewStringWriteable(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, stringWriteable, category)
	settings := &StringWriteableSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	var count = 2
	if settings != nil {
		count = settings.InputCount
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeString, nil, count, 2, 20, body.Inputs, false)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildStringWriteableSchema())
	body.SetDynamicInputs()
	return &StringWriteable{body}, nil
}

func (inst *StringWriteable) Process() {
	count := inst.InputsLen()
	inputs := conversions.ConvertInterfaceToStringMultiple(inst.ReadMultiple(count))
	for _, val := range inputs {
		if val != nil {
			inst.WritePin(node.Out, *val)
			return
		}
	}
	inst.WritePinNull(node.Out)
}

// Custom Node Settings Schema

type StringWriteableSettingsSchema struct {
	InputCount schemas.Integer `json:"inputCount"`
}

type StringWriteableSettings struct {
	InputCount int `json:"inputCount"`
}

func buildStringWriteableSchema() *schemas.Schema {
	props := &StringWriteableSettingsSchema{}
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
