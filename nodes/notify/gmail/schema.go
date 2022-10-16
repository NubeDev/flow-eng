package gmail

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Address schemas.String `json:"address"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Address.Title = "address"
	props.Address.Default = "noreply@nube-io.com"
	schema.Set(props)
	uiSchema := array.Map{
		"address": array.Map{
			"ui:widget": "password",
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "address",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}
