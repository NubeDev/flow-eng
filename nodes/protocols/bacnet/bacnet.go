package bacnet

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

func buildPayload(value string, valueFloat float64) string {
	var v string = value
	if v == "" {
		v = fmt.Sprintf("%f", valueFloat)
	}
	payload := &points.BacnetPayload{Value: v, UUID: helpers.ShortUUID("pay")}
	p, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(p)
}
