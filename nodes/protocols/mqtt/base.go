package broker

const (
	Category   = "mqtt"
	mqttBroker = "mqtt-broker"
	mqttSub    = "mqtt-subscribe"
	mqttPub    = "mqtt-publish"
)

type mqttStore struct {
	parentID string
	payloads []*mqttPayload
}

type mqttPayload struct {
	NodeUUID    string `json:"node_uuid"`
	Topic       string `json:"topic"`
	Payload     string `json:"payload"`
	IsPublisher bool   `json:"is_publisher"`
}

func addUpdatePayload(nodeUUID string, p *mqttStore, newPayload *mqttPayload) (data *mqttStore, found bool) {
	for i, payload := range p.payloads {
		if payload.NodeUUID == nodeUUID {
			p.payloads[i] = newPayload
			found = true
		}
	}
	if !found {
		p.payloads = append(p.payloads, newPayload)
	}
	return p, found
}
