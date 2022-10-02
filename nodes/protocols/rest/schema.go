package rest

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Sch schemas.EnumString `json:"method"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Sch.Title = "method"
	props.Sch.Default = post
	props.Sch.Options = []string{post, patch, put, httpDelete}
	props.Sch.EnumName = []string{post, patch, put, httpDelete}
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

type nodeSettings struct {
	Method string `json:"method"`
}

func getSettings(body map[string]interface{}) (string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return "", err
	}
	if settings != nil {
		return settings.Method, nil
	}
	return "", nil
}
