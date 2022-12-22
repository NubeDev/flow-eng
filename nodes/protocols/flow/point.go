package flow

import (
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
	"time"
)

type FFPoint struct {
	*node.Spec
	topic       string
	lastPayload *covPayload
	lastValue   float64
	lastUpdate  time.Time
}

func NewFFPoint(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPoint, category)
	inputs := node.BuildInputs()
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	lastUpdated := node.BuildOutput(node.LastUpdated, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out, lastUpdated)
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body = node.SetNoParent(body)
	pnt := &FFPoint{body, "", nil, 0, time.Now()}
	return pnt, nil
}

func (inst *FFPoint) set() {
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

func (inst *FFPoint) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		selectedPoint, err := getPointSettings(inst.GetSettings())
		var setTopic bool
		if selectedPoint != nil && err == nil {
			if selectedPoint.Point != "" {
				t := makePointTopic(selectedPoint.Point)
				if t != "" {
					inst.topic = t
					inst.set()
					setTopic = true
				}
			}
		}
		if !setTopic {
			inst.SetWaringMessage("no point selected")
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
			if inst.lastValue != value {
				inst.lastValue = value
				inst.lastUpdate = time.Now()
			} else {
				inst.WritePin(node.LastUpdated, ttime.TimeSince(inst.lastUpdate))
			}
		}
	}
	if !wroteValue {
		inst.WritePinNull(node.Out)
	}
}

func (inst *FFPoint) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
