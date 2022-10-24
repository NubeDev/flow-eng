package flow

import (
	"encoding/json"
)

const (
	category       = "flow"
	flowNetwork    = "flow-network"
	flowPoint      = "flow-point"
	flowPointWrite = "flow-point-write"
)

type covPayload struct {
	Value    float64 `json:"value"`
	ValueRaw float64 `json:"value_raw"`
	Ts       string  `json:"ts"`
	Priority int     `json:"priority"`
}

type PointWriter struct {
	Priority *map[string]*float64 `json:"priority"`
}

// MqttPoint body for getting points from FF over mqtt (can get by name's or uuid, publish on topic rubix/platform/list/points)
type MqttPoint struct {
	NetworkName string       `json:"network_name,omitempty"`
	DeviceName  string       `json:"device_name,omitempty"`
	PointName   string       `json:"point_name,omitempty"`
	PointUUID   string       `json:"point_uuid,omitempty"`
	Priority    *PointWriter `json:"priority,omitempty"`
}

type point struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type pointStore struct {
	parentID string
	points   []*point
	payloads []*pointDetails
}

type pointDetails struct {
	nodeUUID       string
	topic          string
	netDevPntNames string
	pointUUID      string
	payload        string
	isWriteable    bool
}

func parseCOV(body string) (*covPayload, error) {
	marshal, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := &covPayload{}
	err = json.Unmarshal(marshal, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil

}

func getPayloads(children interface{}, ok bool) []*pointDetails {
	if ok {
		mqttData := children.(*pointStore)
		if mqttData != nil {
			return mqttData.payloads
		}
	}
	return nil
}

func addUpdatePayload(nodeUUID string, p *pointStore, newPayload *pointDetails) (data *pointStore, found bool) {
	for i, payload := range p.payloads {
		if payload.nodeUUID == nodeUUID {
			p.payloads[i] = newPayload
			found = true
		}
	}
	if !found {
		p.payloads = append(p.payloads, newPayload)
	}
	return p, found
}
