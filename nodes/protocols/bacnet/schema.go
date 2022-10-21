package bacnet

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type nodeSchema struct {
	Io      schemas.EnumString `json:"ioType"`
	Decimal schemas.Number     `json:"decimal"`
}

func buildSchemaUI() *schemas.Schema {
	props := &nodeSchema{}
	props.Io.Title = "io-type"
	props.Io.Default = string(points.IoTypeVolts)
	props.Io.Options = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent)}
	props.Io.EnumName = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent)}
	props.Decimal.Title = "decimal places"
	props.Decimal.Default = 2
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
	props.Io.Title = "io-type"
	props.Io.Default = string(points.IoTypeVolts)
	props.Io.Options = []string{string(points.IoTypeVolts), string(points.IoTypeDigital)}
	props.Io.EnumName = []string{string(points.IoTypeVolts), string(points.IoTypeDigital)}
	props.Decimal.Title = "decimal places"
	props.Decimal.Default = 2
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
	Io      string `json:"ioType"`
	Decimal int    `json:"decimal"`
}

func getSettings(body map[string]interface{}) (*nodeSettings, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return settings, err
	}
	if settings != nil {
		return settings, nil
	}
	return settings, nil
}
