package bacnetio

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type serverSchema struct {
	AppType     schemas.EnumString `json:"appType"`
	DeviceCount schemas.EnumString `json:"deviceCount"`
	Serial      schemas.EnumString `json:"serial"`
}

var serialPorts = []string{"RC5 485-1", "RC5 485-2", "RC5/RC-IO SIDE-485-PORT", "/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyAMA0"}

func buildSchemaServer() *schemas.Schema {
	props := &serverSchema{}
	props.AppType.Title = "hardware-type"
	props.AppType.Default = string(names.Modbus)
	props.AppType.EnumName = []string{string(names.Modbus), string(names.RubixIO), string(names.RubixIOAndModbus), string(names.Edge)}
	props.AppType.Options = []string{string(names.Modbus), string(names.RubixIO), string(names.RubixIOAndModbus), string(names.Edge)}

	props.DeviceCount.Title = "IO-16-device-count"
	props.DeviceCount.Default = string(names.Modbus)
	props.DeviceCount.EnumName = []string{"1x IO-16", "2x IO-16", "3x IO-16", "4x IO-16"}
	props.DeviceCount.Options = []string{"1", "2", "3", "4"}

	props.Serial.Title = "serial-port (baud-rate:38400)"
	props.Serial.Default = serialPorts[0]
	props.Serial.EnumName = serialPorts
	props.Serial.Options = serialPorts
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "server settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type BacnetSchema struct {
	AppType     string `json:"appType"`
	DeviceCount string `json:"deviceCount"`
	Serial      string `json:"serial"`
}

func getBacnetSchema(body map[string]interface{}) (*BacnetSchema, error) {
	settings := &BacnetSchema{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return settings, err
	}
	if settings != nil {
		return settings, nil
	}
	return settings, nil
}

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
			Title:      "point-settings",
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
			Title:      "point-settings",
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
