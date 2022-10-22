package broker

const (
	category   = "mqtt"
	mqttBroker = "mqtt-broker"
	mqttSub    = "mqtt-subscribe"
	mqttPub    = "mqtt-publish"
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
