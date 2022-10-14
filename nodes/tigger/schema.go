package trigger

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Value schemas.Number `json:"value"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Value.Title = "Set output decimal places"
	props.Value.Default = 2
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "set",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Value int `json:"value"`
}

func getSettings(body map[string]interface{}) (int, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return 0, err
	}
	if settings != nil {
		return settings.Value, nil
	}
	return 0, nil
}
