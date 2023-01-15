package iterator

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Name    schemas.String     `json:"name"`
	Enable  schemas.EnumString `json:"enable"`
	Message schemas.String     `json:"message"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Name.Title = "Name"
	props.Name.Default = "Inject"
	props.Enable.Title = "Enable"
	props.Enable.Options = append(props.Enable.Options, "true", "false")
	props.Enable.Default = "true"
	props.Message.Title = "Message"
	props.Message.Default = "default"
	schema.Set(props)

	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Inject node settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}
