package point

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type NumericWriteable struct {
	*node.Spec
}

func NewNumericWriteable(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, numericWriteable, Category)
	settings := &NumericWriteableSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	var count = 2
	if settings != nil {
		count = settings.InputCount
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, count, 2, 20, body.Inputs, false)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildNumericWriteableSchema())
	body.SetDynamicInputs()
	return &NumericWriteable{body}, nil
}

func (inst *NumericWriteable) Process() {
	// fmt.Println("NumericWriteable() Process()")
	count := inst.InputsLen()
	inputs := conversions.ConvertInterfaceToFloatMultiple(inst.ReadMultiple(count))
	for _, val := range inputs {
		if val != nil {
			// fmt.Println("NumericWriteable() i: ", i, "  val:", *val)
			inst.WritePinFloat(node.Out, *val)
			return
		} else {
			// fmt.Println("NumericWriteable() i: ", i, "  val: null")
		}

	}
	inst.WritePinNull(node.Out)
}

// Custom Node Settings Schema

type NumericWriteableSettingsSchema struct {
	InputCount schemas.Integer `json:"inputCount"`
}

type NumericWriteableSettings struct {
	InputCount int `json:"inputCount"`
}

func buildNumericWriteableSchema() *schemas.Schema {
	props := &NumericWriteableSettingsSchema{}
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
