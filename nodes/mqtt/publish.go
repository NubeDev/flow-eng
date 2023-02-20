package broker

import (
	"github.com/NubeDev/flow-eng/node"
)

type MqttPub struct {
	*node.Spec
	topic string
}

func NewMqttPub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, mqttPub, category)
	top := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs, false, false)
	msg := node.BuildInput(node.Message, node.TypeString, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(top, msg)
	body = node.BuildNode(body, inputs, nil, nil)
	return &MqttPub{body, ""}, nil
}

func (inst *MqttPub) set() {
	s := inst.GetStore()
	parentId := inst.GetParentId()
	nodeUUID := inst.GetID()
	d, ok := s.Get(parentId)
	var mqttData *mqttStore
	if !ok {
		s.Set(parentId, &mqttStore{
			parentID: parentId,
			payloads: []*mqttPayload{&mqttPayload{
				nodeUUID:    nodeUUID,
				topic:       inst.topic,
				isPublisher: true,
			}},
		}, 0)
	} else {
		mqttData = d.(*mqttStore)
		payload := &mqttPayload{
			nodeUUID:    nodeUUID,
			topic:       inst.topic,
			isPublisher: true,
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
