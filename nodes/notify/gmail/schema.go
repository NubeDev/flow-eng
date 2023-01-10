package gmail

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	FromAddress schemas.String `json:"fromAddress"`
	Password    schemas.String `json:"password"`
	ToAddress   schemas.String `json:"toAddress"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.FromAddress.Title = "From Address"
	props.FromAddress.Default = "noreply@nube-io.com"
	props.ToAddress.Title = "To Address"
	schema.Set(props)
	uiSchema := array.Map{
		"password": array.Map{
			"ui:widget": "password",
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "addresses and password",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}
