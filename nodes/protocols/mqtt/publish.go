package broker

import (
	"github.com/NubeDev/flow-eng/node"
)

type MqttPub struct {
	*node.Spec
	topic string
}

const pubHelp = `A node for publishing an MQTT message to a broker. (must be added inside the MQTT broker node)`

func NewMqttPub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, mqttPub, category)
	top := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp("mqtt topic, leave blank if the topic comes in the message body"))
	msg := node.BuildInput(node.Message, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp("body to send in the mqtt-message, if you wanna send json use msg.body{} and you can also send the topic as msg.topic=myTopic"))
	inputs := node.BuildInputs(top, msg)
	body = node.BuildNode(body, inputs, nil, nil)
	body.SetHelp(pubHelp)
	return &MqttPub{body, ""}, nil
}

func (inst *MqttPub) set() {
	s := inst.GetStore()
	parentId := inst.GetParentId()
	nodeUUID := inst.GetID()
	d, ok := s.Get(parentId)
	var mqttData *mqttStore
	if !ok { // never added so add
		s.Set(parentId, &mqttStore{
			parentID: parentId,
			payloads: []*mqttPayload{&mqttPayload{
				NodeUUID:    nodeUUID,
				Topic:       inst.topic,
				IsPublisher: true,
			}},
		}, 0)
	} else { // add new payload to existing
		mqttData = d.(*mqttStore)
		payload := &mqttPayload{
			NodeUUID:    nodeUUID,
			Topic:       inst.topic,
			IsPublisher: true,
		}
		mqttData, _ = addUpdatePayload(nodeUUID, mqttData, payload)
		s.Set(parentId, mqttData, 0)
	}
}

func (inst *MqttPub) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		topic, null := inst.ReadPinAsString(node.Topic)
		if !null {
			inst.topic = topic
			inst.set()
		}
	}

}
