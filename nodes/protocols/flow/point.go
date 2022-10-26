package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
)

type Point struct {
	*node.Spec
	topic       string
	lastPayload *covPayload
}

func NewPoint(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPoint, category)
	inputs := node.BuildInputs()
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	pnt := &Point{body, "", nil}
	body.SetSchema(pnt.buildSchema())
	return pnt, nil
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
		fmt.Println("SETTINGS", topic.Point, inst.GetSettings(), err)
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
	val, null := inst.GetPayloadNull()
	if null {
		inst.WritePinNull(node.Out)
	} else {
		p, err := parseCOV(val)
		if err == nil && p != nil {
			inst.lastPayload = p
			inst.WritePinFloat(node.Out, p.Value, 2)
		} else {
			inst.WritePinNull(node.Out)
		}
	}
}

func (inst *Point) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
