package trigger

import (
	"github.com/NubeDev/flow-eng/nodes/trigger"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Name         schemas.String     `json:name`
	Interval     schemas.Number     `json:interval`
	Units        schemas.EnumString `json:units`
	COVThreshold schemas.Number     `json:"covThreshold"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.COVThreshold.Title = "COV Threshold"
	props.COVThreshold.Default = 2
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Set output decimal places",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Name         string            `json:name`
	Interval     float64           `json:interval`
	Units        trigger.TimeUnits `json:units`
	COVThreshold int               `json:"covThreshold"`
}

func getSettings(body map[string]interface{}) (int, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return 0, err
	}
	if settings != nil {
		return settings.COVThreshold, nil
	}
	return 0, nil
}
