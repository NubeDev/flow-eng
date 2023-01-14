package cov

import (
	"github.com/NubeDev/flow-eng/nodes/trigger"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Name         schemas.String     `json:"name"`
	Interval     schemas.Number     `json:"interval"`
	Units        schemas.EnumString `json:"units"`
	COVThreshold schemas.Number     `json:"covThreshold"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Name.Title = "Name"
	props.Name.Default = "COV"
	props.Interval.Title = "COV Interval"
	props.Interval.Default = 2
	props.Units.Title = "Time units"
	props.Units.Default = string(trigger.Seconds)
	props.Units.Options = append(props.Units.Options, string(trigger.Milliseconds), string(trigger.Seconds), string(trigger.Minutes), string(trigger.Hours))
	props.COVThreshold.Title = "COV Threshold"
	props.COVThreshold.Default = 5
	schema.Set(props)

	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "COV settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}
