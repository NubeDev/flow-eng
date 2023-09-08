package bacnetio

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type serverSchema struct {
	DeviceCount    schemas.EnumString   `json:"deviceCount"`
	Serial         schemas.EnumString   `json:"serial"`
	Timeout        schemas.NumberLimits `json:"timeout"`
	DeviceAddress1 schemas.NumberLimits `json:"deviceAddress1"`
	DeviceAddress2 schemas.NumberLimits `json:"deviceAddress2"`
	DeviceAddress3 schemas.NumberLimits `json:"deviceAddress3"`
	DeviceAddress4 schemas.NumberLimits `json:"deviceAddress4"`
}

const noDevices = "Select an IO-16"

var ioOptions = []string{noDevices, "1x IO-16", "2x IO-16s", "3x IO-16s", "4x IO-16s"}
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
	props.Timeout.Title = "timeout (ms)"
	props.Timeout.Default = 100
	props.Timeout.Min = 100
	props.Timeout.Max = 20000

	// device address
	var min int = 1
	var max int = 50

	props.DeviceAddress1.Title = "address device-1"
	props.DeviceAddress1.Default = 1
	props.DeviceAddress1.Min = min
	props.DeviceAddress1.Max = max

	props.DeviceAddress2.Title = "address device-2"
	props.DeviceAddress2.Default = 2
	props.DeviceAddress2.Min = min
	props.DeviceAddress2.Max = max

	props.DeviceAddress3.Title = "address device-3"
	props.DeviceAddress3.Default = 3
	props.DeviceAddress3.Min = min
	props.DeviceAddress3.Max = max

	props.DeviceAddress4.Title = "address device-4"
	props.DeviceAddress4.Default = 4
	props.DeviceAddress4.Min = min
	props.DeviceAddress4.Max = max

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
	DeviceCount    string `json:"deviceCount"`
	Serial         string `json:"serial"`
	Timeout        int    `json:"timeout"`
	DeviceAddress1 int    `json:"deviceAddress1"`
	DeviceAddress2 int    `json:"deviceAddress2"`
	DeviceAddress3 int    `json:"deviceAddress3"`
	DeviceAddress4 int    `json:"deviceAddress4"`
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
