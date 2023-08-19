package rql

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Rule    schemas.String     `json:"name"`
	Results schemas.EnumString `json:"results"`
}

const (
	resultLatest = "latest result"
	resultOldest = "oldest result"
	resultAll    = "all results"
)

func (inst *Get) buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Rule.Title = "Enter rule by name"

	props.Results.Title = "selected results"
	props.Results.Default = resultLatest
	opts := []string{resultLatest, resultOldest, resultAll}
	props.Results.Options = opts
	props.Results.EnumName = opts

	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "RQL get rule",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Name    string `json:"name"`
	Results string `json:"results"`
}

func (inst *Get) getSettings() (*nodeSettings, error) {
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
