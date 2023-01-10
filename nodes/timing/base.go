package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
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
	Time     schemas.EnumString `json:"time"`
	Duration schemas.Number     `json:"duration"`
}

type defaultNodeSettings struct {
	Time     string        `json:"time"`
	Duration time.Duration `json:"duration"`
}

func buildDefaultSchema() *schemas.Schema {
	props := &defaultNodeSchema{}
	// time selection
	props.Duration.Title = "duration"
	props.Duration.Default = 1

	// time selection
	props.Time.Title = "time"
	props.Time.Default = ttime.Sec
	props.Time.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.Time.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	pprint.PrintJSON(props)
	schema.Set(props)

	fmt.Println(fmt.Sprintf("buildSchema() props: %+v", props))
	pprint.PrintJSON(props)

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
	fmt.Println(fmt.Sprintf("buildSchema() s: %+v", s))
	pprint.PrintJSON(s)
	return s
}

func getSettings(body map[string]interface{}) (*defaultNodeSettings, error) {
	settings := &defaultNodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
