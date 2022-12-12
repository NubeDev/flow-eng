package bacnetio

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type serverSchema struct {
	DeviceCount schemas.EnumString `json:"deviceCount"`
	Serial      schemas.EnumString `json:"serial"`
}

var serialPorts = []string{"485-1", "485-2", "SIDE-485-PORT", "/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyAMA0"}

func BuildSchemaServer() *schemas.Schema {
	props := &serverSchema{}
	props.DeviceCount.Title = "IO-16-device-count"
	props.DeviceCount.Default = "No IO-16s"
	props.DeviceCount.EnumName = []string{"No IO-16s", "1x IO-16", "2x IO-16s", "3x IO-16s", "4x IO-16s"}
	props.DeviceCount.Options = []string{"0", "1", "2", "3", "4"}

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
	DeviceCount string `json:"deviceCount"`
	Serial      string `json:"serial"`
}

func GetBacnetSchema(body map[string]interface{}) (*BacnetSchema, error) {
	settings := &BacnetSchema{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
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
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
