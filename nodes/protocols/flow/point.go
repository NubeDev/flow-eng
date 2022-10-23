package flow

import (
	"github.com/NubeDev/flow-eng/node"
)

type Point struct {
	*node.Spec
	topic string
}

func NewPoint(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPoint, category)
	inputs := node.BuildInputs()
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Point{body, ""}, nil
}

func (inst *Point) set() {
	s := inst.GetStore()
	parentId := inst.GetParentId()
	nodeUUID := inst.GetID()
	d, ok := s.Get(parentId)
	var mqttData *pointStore
	if !ok {
		s.Set(parentId, &pointStore{
			parentID: parentId,
			payloads: []*pointDetails{&pointDetails{
				nodeUUID: nodeUUID,
				topic:    inst.topic,
			}},
		}, 0)
	} else {
		mqttData = d.(*pointStore)
		payload := &pointDetails{
			nodeUUID: nodeUUID,
			topic:    inst.topic,
		}
		mqttData, _ = addUpdatePayload(nodeUUID, mqttData, payload)
		s.Set(parentId, mqttData, 0)
	}
}

func (inst *Point) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		topic, err := getPointSettings(inst.GetSettings())
		if err == nil {
			if topic.Point != "" {
				t := pointTopic(topic.Point)
				if t != "" {
					inst.topic = t
					inst.set()
				}
			}
		}
	}
	val, null := inst.ReadPayloadAsString()
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePin(node.Out, val)
	}
}
