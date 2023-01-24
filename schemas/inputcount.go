package schemas

import "github.com/NubeIO/lib-schema/schema"

type InputCount struct {
	NumberLimits NumberLimits `json:"inputCount"`
}

func GetInputCount() *Schema {
	props := &InputCount{}
	props.NumberLimits.Title = "Input Count"
	props.NumberLimits.Default = 2

	props.NumberLimits.Max = 20
	schema.Set(props)

	s := &Schema{
		Schema: SchemaBody{
			Title:      "Input Count",
			Properties: props,
		},
		UiSchema: nil,
	}

	return s
}
