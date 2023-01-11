package timing

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

const (
	category  = "time"
	delay     = "delay"
	delayOn   = "delay-on"
	delayOff  = "delay-off"
	dutyCycle = "duty-cycle"
	minOnOff  = "min-on-off"
	oneShot   = "one-shot"
)

type defaultNodeSchema struct {
	Interval  schemas.Number     `json:"interval"`
	TimeUnits schemas.EnumString `json:"time"`
}

type defaultNodeSettings struct {
	Interval  float64 `json:"interval"`
	TimeUnits string  `json:"time_units"`
}

func buildDefaultSchema() *schemas.Schema {
	props := &defaultNodeSchema{}
	// time selection
	props.Interval.Title = "Interval"
	props.Interval.Default = 1

	// time selection
	props.TimeUnits.Title = "Time Units"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.TimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	schema.Set(props)

	uiSchema := array.Map{
		"time": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "time_units"},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Node Settings",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

func getDefaultSettings(body map[string]interface{}) (*defaultNodeSettings, error) {
	settings := &defaultNodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
