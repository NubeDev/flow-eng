package store

import (
	"fmt"
	"testing"
)

type mqttStore struct {
	brokerUUID string
	payloads   []*mqttPayload
}

type mqttPayload struct {
	nodeUUID string
	topic    string
	payload  string
}

func TestInit(t *testing.T) {
	s := Init()

	s.Set("mqtt", &mqttStore{
		brokerUUID: "123",
		payloads: []*mqttPayload{&mqttPayload{
			nodeUUID: "n123",
			topic:    "abc",
			payload:  "payload",
		}},
	}, 0)

	d, ok := s.Get("mqtt")
	if !ok {
		return
	}
	mqttData := d.(*mqttStore)
	for _, payload := range mqttData.payloads {
		if payload.nodeUUID == "n123" {
			fmt.Println("boo")
		}
	}

}
