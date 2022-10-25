package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
)

type PointWrite struct {
	*node.Spec
	topic string
}

func NewPointWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPointWrite, category)
	in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs)
	in16 := node.BuildInput(node.In16, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(in15, in16)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &PointWrite{body, ""}, nil
}

func (inst *PointWrite) set() {
	s := inst.GetStore()
	parentId := inst.GetParentId()
	nodeUUID := inst.GetID()
	d, ok := s.Get(parentId)
	var mqttData *pointStore
	if !ok {
		s.Set(parentId, &pointStore{
			parentID: parentId,
			payloads: []*pointDetails{&pointDetails{
				nodeUUID:       nodeUUID,
				netDevPntNames: inst.topic,
				isWriteable:    true,
			}},
		}, 0)
	} else {
		mqttData = d.(*pointStore)
		payload := &pointDetails{
			nodeUUID:       nodeUUID,
			netDevPntNames: inst.topic,
			isWriteable:    true,
		}
		mqttData, _ = addUpdatePayload(nodeUUID, mqttData, payload)
		s.Set(parentId, mqttData, 0)
	}
}

func (inst *PointWrite) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		topic, err := getPointSettings(inst.GetSettings())
		fmt.Println("point write", topic.Point)
		if err == nil {
			if topic.Point != "" {
				if topic.Point != "" {
					inst.topic = topic.Point
					inst.set()
				}
			}
		}
	}
	val, null := inst.ReadPayloadAsString()
	if null {
		inst.WritePinNull(node.Out)
	} else {
		p, err := parseCOV(val)
		fmt.Println(p, err)
		inst.WritePin(node.Out, val)
	}
}

func (inst *PointWrite) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
