package cov

import (
	"github.com/NubeDev/flow-eng/nodes/trigger"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Name         schemas.String     `json:"name"`
	Interval     schemas.Number     `json:"interval"`
	Units        schemas.EnumString `json:"units"`
	COVThreshold schemas.Number     `json:"covThreshold"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Name.Title = "Name"
	props.Name.Default = "COV"
	props.Interval.Title = "COV Threshold"
	props.Interval.Default = 2
	props.Units.Title = "Time units"
	props.Units.Default = string(trigger.Seconds)
	props.COVThreshold.Title = "COV Threshold"
	props.COVThreshold.Default = 2
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "COV settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Name         string            `json:"name"`
	Interval     float64           `json:"interval"`
	Units        trigger.TimeUnits `json:"units"`
	COVThreshold int               `json:"covThreshold"`
}

func getSettings(body map[string]interface{}) (map[string]interface{}, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, &settings)
	r := make(map[string]interface{}, 4)
	if err != nil {
		return r, err
	}
	if settings != nil {
		r["Name"] = settings.Name
		r["Interval"] = settings.Interval
		r["Units"] = settings.Units
		r["COVThreshold"] = settings.COVThreshold
		return r, nil
	}
	return r, nil
}
