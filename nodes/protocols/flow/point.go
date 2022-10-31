package flow

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
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
	body = node.SetNoParent(body)
	pnt := &Point{body, "", nil}
	return pnt, nil
}

func (inst *Point) set() {
	s := inst.GetStore()
	if s == nil {
		return
	}
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
		selectedPoint, err := getPointSettings(inst.GetSettings())
		var setTopic bool
		if selectedPoint != nil && err == nil {
			if selectedPoint.Point != "" {
				if selectedPoint.Point != "" {
					t := makePointTopic(selectedPoint.Point)
					if t != "" {
						inst.topic = t
						inst.set()
					}
				}
			}
		}
		if !setTopic {
			inst.SetWaringMessage("no point has been selected")
			inst.SetWaringIcon(string(emoji.OrangeCircle))
		}
	}
	val, null := inst.GetPayloadNull()
	var wroteValue bool
	if null {
		inst.WritePinNull(node.Out)
	} else {
		p, value, _, err := parseCOV(val)
		if err == nil && p != nil {
			inst.lastPayload = p
			wroteValue = true
			inst.WritePinFloat(node.Out, value, 2)
		}
	}
	if !wroteValue {
		inst.WritePinNull(node.Out)
	}
}

func (inst *Point) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
