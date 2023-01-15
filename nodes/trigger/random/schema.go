package random

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Name      schemas.String `json:"name"`
	Percision schemas.Number `json:"percision"`
	Max       schemas.Number `json:"max"`
	Min       schemas.Number `json:"min"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Name.Title = "Name"
	props.Name.Default = "Random"
	props.Percision.Title = "Percision"
	props.Percision.Default = 2
	props.Max.Title = "Max"
	props.Max.Default = 5
	props.Min.Title = "Min"
	props.Min.Default = 0
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Random node settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

// type nodeSettings struct {
// 	Name         string            `json:name`
// 	Interval     float64           `json:interval`
// 	Units        trigger.TimeUnits `json:units`
// 	COVThreshold int               `json:"covThreshold"`
// }

// func getSettings(body map[string]interface{}) (int, error) {
// 	settings := &nodeSettings{}
// 	err := mapstructure.Decode(body, settings)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if settings != nil {
// 		return settings.COVThreshold, nil
// 	}
// 	return 0, nil
// }
