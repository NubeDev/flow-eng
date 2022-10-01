package rest

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Address schemas.EnumString `json:"address"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Address.Title = "method"
	props.Address.Default = get
	props.Address.Options = []string{post, patch, put, httpDelete}
	props.Address.EnumName = []string{post, patch, put, httpDelete}
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "method",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}
