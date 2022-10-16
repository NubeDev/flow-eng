package timing

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Time     schemas.EnumString `json:"time"`
	Duration schemas.Number     `json:"function"`
}

// NODEs will single in/out
const (
	ms  = "ms"
	sec = "sec"
	min = "min"
	hr  = "hour"
)

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	// time selection
	props.Duration.Title = "duration"
	props.Duration.Default = 1

	// time selection
	props.Time.Title = "time"
	props.Time.Default = sec
	props.Time.Options = []string{ms, sec, min, hr}
	props.Time.EnumName = []string{ms, sec, min, hr}

	schema.Set(props)
	uiSchema := array.Map{
		"time": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Set delay time",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

type nodeSettings struct {
	Function string `json:"function"`
}

func getSettings(body map[string]interface{}) (string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return "", err
	}
	if settings != nil {
		return settings.Function, nil
	}
	return "", nil
}
