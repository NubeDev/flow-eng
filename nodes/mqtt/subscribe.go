package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/services/mqttclient"
)

type MqttSub struct {
	*node.Spec
	client     *mqttclient.Client
	connected  bool
	subscribed bool
	topic      string
}

const (
	topic = "topic"
)

func NewMqttSub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, mqttSub, category)
	top := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	ip := node.BuildInput(node.Ip, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(top, ip)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &MqttSub{body, nil, false, false, ""}, nil
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
	topic, _ := inst.ReadPinAsString(node.Topic)
	inst.topic = topic
	inst.set()
	v, null := inst.ReadPayloadAsString()
	fmt.Println("MQTT-SUB Process()", v, null)
}
