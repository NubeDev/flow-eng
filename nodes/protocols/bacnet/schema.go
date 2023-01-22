package bacnetio

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type serverSchema struct {
	DeviceCount schemas.EnumString `json:"deviceCount"`
	Serial      schemas.EnumString `json:"serial"`
}

func selectIO16() string {
	for i, _ := range ioOptions {
		fmt.Println("!!!!!!!!!1", i)
		return fmt.Sprintf("%d", i)
	}
	return "0"
}

var ioOptions = []string{"Select an IO-16", "1x IO-16", "2x IO-16s", "3x IO-16s", "4x IO-16s"}
var serialPorts = []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyAMA0"}

func BuildSchemaServer() *schemas.Schema {
	props := &serverSchema{}
	props.DeviceCount.Title = "IO-16-device-count"
	props.DeviceCount.Default = ioOptions[1]
	props.DeviceCount.EnumName = ioOptions
	props.DeviceCount.Options = ioOptions
	props.Serial.Title = "serial-port (baud-rate:38400)"
	props.Serial.Default = serialPorts[2]
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
	Offset  schemas.Number     `json:"offset"`
}

func buildSchemaUI() *schemas.Schema {
	props := &nodeSchema{}
	props.Io.Title = "io-type"
	props.Io.Default = string(points.IoTypeVolts)
	props.Io.Options = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent)}
	props.Io.EnumName = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent)}

	props.Decimal.Title = "decimal places"
	props.Decimal.Default = 2

	props.Offset.Title = "offset"
	props.Offset.Default = 0
	props.Offset.Minimum = -100000000
	props.Offset.Maximum = 100000000
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

	props.Offset.Title = "offset"
	props.Offset.Default = 0
	props.Offset.Minimum = -100000000
	props.Offset.Maximum = 100000000

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
	Io      string  `json:"ioType"`
	Decimal int     `json:"decimal"`
	Offset  float64 `json:"offset"`
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

func bacnetAddress(deviceCount int, prefix, ioTypePrefix string) []string {
	var devicesCount = 1
	var ioCount = 1
	var out []string
	if deviceCount == 0 {
		deviceCount = 1
	}
	var devCount = deviceCount * 8
	for i := 0; i <= devCount; i++ {
		if i > 0 {
			if i%8 == 0 {
				devicesCount++
			}
		}
		if i%8 == 0 {
			ioCount = 1
		}
		dev := fmt.Sprintf("dev-%d", devicesCount)
		bacnet := fmt.Sprintf("%s-%d", prefix, i+1)
		value := fmt.Sprintf("%s %s %s-%d", dev, bacnet, ioTypePrefix, ioCount)
		out = append(out, value)
		ioCount++
	}
	return out
}
