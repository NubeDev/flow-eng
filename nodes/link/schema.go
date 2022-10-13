package link

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Sch schemas.EnumString `json:"topic"`
}

func options() []string {
	s := getStore()
	var out []string
	for _, v := range s.GetAll() {
		out = append(out, v.topic)
	}
	return out
}

func buildSchema() *schemas.Schema {
	opts := options()
	props := &nodeSchema{}
	props.Sch.Title = "Select Topic"
	props.Sch.Default = ""
	props.Sch.Options = opts
	props.Sch.EnumName = opts
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Select Topic",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Topic string `json:"topic"`
}

func getSettings(body map[string]interface{}) (string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return "", err
	}
	if settings != nil {
		return settings.Topic, nil
	}
	return "", nil
}
