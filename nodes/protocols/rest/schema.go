package rest

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Sch schemas.EnumString `json:"method"`
}

func (inst *HTTP) buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Sch.Title = "method"
	props.Sch.Default = get
	props.Sch.Options = []string{get, post, patch, put, httpDelete}
	props.Sch.EnumName = []string{get, post, patch, put, httpDelete}
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "method",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Method string `json:"method"`
}

func (inst *HTTP) getSettings() (*nodeSettings, error) {
	body := inst.GetSettings()
	settings := &nodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	if settings == nil {
		return nil, errors.New("settings are empty")
	}
	return settings, nil
}
