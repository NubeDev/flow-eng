package timing

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type nodeSchema struct {
	Time     schemas.EnumString `json:"time"`
	Duration schemas.Number     `json:"duration"`
}

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
	Time     string        `json:"time"`
	Duration time.Duration `json:"duration"`
}

func getSettings(body map[string]interface{}) (*nodeSettings, error) {
	settings := &nodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
