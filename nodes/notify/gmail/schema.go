package gmail

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	FromAddress schemas.String `json:"fromAddress"`
	Token       schemas.String `json:"token"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.FromAddress.Title = "From Address"
	props.FromAddress.Default = "noreply@nube-io.com"
	schema.Set(props)
	uiSchema := array.Map{
		"token": array.Map{
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
