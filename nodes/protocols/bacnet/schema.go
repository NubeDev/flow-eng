package bacnet

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Sch schemas.EnumString `json:"ioType"`
}

func buildSchemaUI() *schemas.Schema {
	props := &nodeSchema{}
	props.Sch.Title = "io-type"
	props.Sch.Default = string(points.IoTypeVolts)
	props.Sch.Options = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent)}
	props.Sch.EnumName = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent)}
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "io-type",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

func buildSchemaUO() *schemas.Schema {
	props := &nodeSchema{}
	props.Sch.Title = "io-type"
	props.Sch.Default = string(points.IoTypeVolts)
	props.Sch.Options = []string{string(points.IoTypeVolts), string(points.IoTypeDigital)}
	props.Sch.EnumName = []string{string(points.IoTypeVolts), string(points.IoTypeDigital)}
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "io-type",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	IoType string `json:"ioType"`
}

func getSettings(body map[string]interface{}) (string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return "", err
	}
	if settings != nil {
		return settings.IoType, nil
	}
	return "", nil
}
