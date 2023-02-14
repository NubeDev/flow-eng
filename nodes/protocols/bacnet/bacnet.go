package bacnetio

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

func buildPayload(value string, valueFloat float64) string {
	if value == "" {
		value = fmt.Sprintf("%f", valueFloat)
	}
	payload := &points.BacnetPayload{Value: value, UUID: helpers.ShortUUID("pay")}
	p, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(p)
}

func buildPayloadName(value string) string {
	payload := &points.BacnetPayload{Value: value, UUID: helpers.ShortUUID("pay")}
	p, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(p)
}

func buildPayloadWithoutUUID(value string, valueFloat float64) string {
	if value == "" {
		value = fmt.Sprintf("%f", valueFloat)
	}
	payload := &points.BacnetPayload{Value: value, UUID: ""}
	p, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(p)
}
