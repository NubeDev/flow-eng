package iterator

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Name       schemas.String     `json:"name"`
	Interval   schemas.Number     `json:"interval"`
	Units      schemas.EnumString `json:"units"`
	Iterations schemas.Number     `json:"iterations"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Name.Title = "Name"
	props.Name.Default = "Iterator"

	props.Interval.Title = "Interval"
	props.Interval.Default = 4

	props.Units.Title = "Units"
	props.Units.Options = append(props.Units.Options, Milliseconds, Seconds, Minutes, Hours)
	props.Units.Default = Seconds

	props.Iterations.Title = "Iterations"
	props.Iterations.Default = 10
	schema.Set(props)

	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Iterator node settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}
