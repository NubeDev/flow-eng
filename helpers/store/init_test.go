package store

import (
	pprint "github.com/NubeDev/flow-eng/helpers/pprint"
	"testing"
)

type mqttStore struct {
	ParentId string         `json:"parent_id"`
	Payloads []*mqttPayload `json:"payloads"`
}

type mqttPayload struct {
	NodeUUID string `json:"node_uuid"`
	Topic    string `json:"topic"`
	Payload  string `json:"payload"`
}

func addUpdatePayload(nodeUUID string, p *mqttStore, newPayload *mqttPayload) (data *mqttStore, found bool) {
	for i, payload := range p.Payloads {
		if payload.NodeUUID == nodeUUID {
			p.Payloads[i] = newPayload
			found = true
		}
	}
	return p, found
}

func TestInit(t *testing.T) {
	s := Init()
	var parentId = "123"
	var nodeUUID = "n123"

	// try and get if nil then set
	// if not nil the get parentId
	// on the parentId set new payload for the subscriber

	// set node n123
	s.Set("mqtt", &mqttStore{
		ParentId: parentId,
		Payloads: []*mqttPayload{&mqttPayload{
			NodeUUID: nodeUUID,
			Topic:    "abc",
			Payload:  "payload",
		}},
	}, 0)

	d, ok := s.Get("mqtt")
	if !ok {
		return
	}
	mqttData := d.(*mqttStore)

	mqttData, _ = addUpdatePayload(nodeUUID, mqttData, &mqttPayload{
		NodeUUID: nodeUUID,
		Topic:    "xyz",
		Payload:  "payload2",
	})
	pprint.PrintJSON(mqttData)
	d, ok = s.Get("mqtt")
	if !ok {
		return
	}
	mqttData = d.(*mqttStore)

}
