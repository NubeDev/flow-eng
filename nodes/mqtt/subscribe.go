package broker

import (
	"github.com/NubeDev/flow-eng/node"
)

type MqttSub struct {
	*node.Spec
	topic string
}

func NewMqttSub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, mqttSub, category)
	top := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(top)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &MqttSub{body, ""}, nil
}

func (inst *MqttSub) set() {
	s := inst.GetStore()
	parentId := inst.GetParentId()
	nodeUUID := inst.GetID()
	d, ok := s.Get(parentId)
	var mqttData *mqttStore
	if !ok {
		s.Set(parentId, &mqttStore{
			parentID: parentId,
			payloads: []*mqttPayload{&mqttPayload{
				nodeUUID: nodeUUID,
				topic:    inst.topic,
			}},
		}, 0)
	} else {
		mqttData = d.(*mqttStore)
		payload := &mqttPayload{
			nodeUUID: nodeUUID,
			topic:    inst.topic,
		}
		mqttData, _ = addUpdatePayload(nodeUUID, mqttData, payload)
		s.Set(parentId, mqttData, 0)
	}
}

func (inst *MqttSub) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		topic, null := inst.ReadPinAsString(node.Topic)
		if !null {
			inst.topic = topic
			inst.set()
		}
	}
	val, null := inst.ReadMQTTPayloadAsString()
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePin(node.Out, val)
	}
}
