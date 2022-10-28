package flow

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
)

type PointWrite struct {
	*node.Spec
	topic          string
	netDevicePoint string
}

func NewPointWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPointWrite, category)
	in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs)
	in16 := node.BuildInput(node.In16, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(in15, in16)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body = node.SetNoParent(body)
	return &PointWrite{body, "", ""}, nil
}

func (inst *PointWrite) set() {
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
				nodeUUID:       nodeUUID,
				topic:          inst.topic,
				netDevPntNames: inst.netDevicePoint,
				isWriteable:    true,
			}},
		}, 0)
	} else {
		mqttData = d.(*pointStore)
		payload := &pointDetails{
			nodeUUID:       nodeUUID,
			topic:          inst.topic,
			netDevPntNames: inst.netDevicePoint,
			isWriteable:    true,
		}
		mqttData, _ = addUpdatePayload(nodeUUID, mqttData, payload)
		s.Set(parentId, mqttData, 0)
	}
}

func (inst *PointWrite) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		selectedPoint, err := getPointSettings(inst.GetSettings())
		var setTopic bool
		if selectedPoint != nil && err == nil {

			if selectedPoint.Point != "" {
				t := makePointTopic(selectedPoint.Point)
				if t != "" {
					inst.topic = t
					inst.netDevicePoint = selectedPoint.Point
					inst.set()
					setTopic = true
				}
			}
		}
		if !setTopic {
			inst.SetWaringMessage("no point has been selected")
		}
	}

	val, null := inst.GetPayloadNull()
	var wroteValue bool
	if null {
		inst.WritePinNull(node.Out)
	} else {
		_, value, _, err := parseCOV(val)
		if err == nil {
			wroteValue = true
			inst.WritePin(node.Out, value)
		}
	}
	if !wroteValue {
		inst.WritePinNull(node.Out)
	}
}

func (inst *PointWrite) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
