package broker

import "github.com/NubeDev/flow-eng/node"

const subHelp = `A node for subscribing to an MQTT topic and message to a broker. (must be added inside the MQTT broker node)`

type MqttSub struct {
	*node.Spec
	topic string
}

func NewMqttSub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, mqttSub, category)
	top := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp("mqtt topic example: my/topic"))
	inputs := node.BuildInputs(top)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs, node.SetOutputHelp(node.OutHelp)))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(subHelp)
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
				NodeUUID: nodeUUID,
				Topic:    inst.topic,
			}},
		}, 0)
	} else {
		mqttData = d.(*mqttStore)
		payload := &mqttPayload{
			NodeUUID: nodeUUID,
			Topic:    inst.topic,
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
