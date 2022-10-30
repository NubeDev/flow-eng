package link

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type nodeSchema struct {
	Sch schemas.EnumString `json:"topic"`
}

func options(nodeType string) []string {
	s := getStore()
	var out []string
	out = append(out, "select a topic")
	for _, v := range s.GetAll() {
		parts := strings.Split(v.topic, "-")
		if len(parts) > 0 {
			if parts[0] == nodeType {
				out = append(out, v.topic)
			}
		}

	}
	return out
}

func buildSchema(nodeType string) *schemas.Schema {
	opts := options(nodeType)
	props := &nodeSchema{}
	props.Sch.Title = "Select Topic"
	if len(opts) > 0 {
		props.Sch.Default = opts[0]
	}
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
