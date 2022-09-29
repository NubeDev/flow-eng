package schemas

import "github.com/NubeIO/lib-schema/schema"

type InputCount struct {
	NumberLimits NumberLimits `json:"inputCount"`
}

func GetInputCount() *Schema {
	m := &InputCount{}
	m.NumberLimits.Title = "count"
	m.NumberLimits.Min = 2
	m.NumberLimits.Max = 20
	schema.Set(m)

	s := &Schema{
		Title:      "Input Count",
		Properties: m,
		UiSchema:   nil,
	}

	return s
}
