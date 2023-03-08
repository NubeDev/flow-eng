package flow

import (
	"fmt"
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
	currentPriority := node.BuildOutput(node.CurrentPriority, node.TypeFloat, nil, body.Outputs)
	lastUpdated := node.BuildOutput(node.LastUpdated, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out, currentPriority, lastUpdated)
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body = node.SetNoParent(body)
	pnt := &FFPoint{body, "", nil, 0, time.Now()}
	return pnt, nil
}

// set add this point to the store
func (inst *FFPoint) set() {
	// log.Infof("FLOW POINT: set() topic: %+v", inst.topic)
	// log.Infof("FLOW POINT: set() STORE: %+v", inst.GetStore().All())
	// log.Infof("FLOW POINT: set() STORE Object: %+v", inst.GetStore().All()[inst.GetParentId()].Object)
	s := inst.GetStore()
	if s == nil {
		return
	}
	parentId := inst.GetParentId()
	// log.Infof("FLOW POINT: set() inst.GetParentId(): %+v", inst.GetParentId())
	nodeUUID := inst.GetID()
	// log.Infof("FLOW POINT: set() nodeUUID: %+v", nodeUUID)
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
		// log.Infof("FLOW POINT: set()")
	} else {
		mqttData = d.(*pointStore)
		/*
			log.Infof("FLOW POINT: set() mqttData: %+v", *mqttData)
			for i, v := range mqttData.payloads {
				if v == nil {
					log.Infof("FLOW POINT: set() mqttData.payloads %v IS NIL", i)
				} else {
					log.Infof("FLOW POINT: set() mqttData.payloads %v: %+v", i, *v)
				}
			}
		*/
		payload := &pointDetails{
			nodeUUID: nodeUUID,
			topic:    inst.topic,
		}
		mqttData, _ = addUpdatePayload(nodeUUID, mqttData, payload)
		s.Set(parentId, mqttData, 0)
		/*
			for i, v := range mqttData.payloads {
				if v == nil {
					log.Infof("FLOW POINT: set() mqttData.payloads %v IS NIL", i)
				} else {
					log.Infof("FLOW POINT: set() mqttData.payloads %v: %+v", i, *v)
				}
			}
		*/
	}
}

func (inst *FFPoint) checkStillExists() {
	s := inst.GetStore()
	if s == nil {
		return
	}
	parentId := inst.GetParentId()
	topic := fmt.Sprintf("pointsList_%s", parentId)
	children, ok := s.Get(topic)
	if ok {
		existingPoints := children.([]*point)
		var pointExists bool
		for _, existingPoint := range existingPoints {
			t := makePointTopic(existingPoint.Name)
			if t == inst.topic {
				pointExists = true
				inst.SetSubTitle(existingPoint.Name)
				inst.SetWaringMessage("")
				inst.SetWaringIcon(string(emoji.GreenCircle))
			}
		}
		if !pointExists {
			inst.SetWaringMessage(pointError)
			inst.SetWaringIcon(string(emoji.OrangeCircle))
			inst.SetSubTitle("")
		}
	}
}

func (inst *FFPoint) setTopic() {
	selectedPoint, err := getPointSettings(inst.GetSettings())
	if selectedPoint != nil && err == nil {
		if selectedPoint.Point != "" {
			t := makePointTopic(selectedPoint.Point)
			if t != "" {
				inst.topic = t
				inst.set()
				inst.SetSubTitle(selectedPoint.Point)
				inst.SetWaringMessage("")
				inst.SetWaringIcon(string(emoji.GreenCircle))
			} else {
				inst.SetWaringMessage(pointError)
				inst.SetWaringIcon(string(emoji.OrangeCircle))
				inst.SetSubTitle("")
			}
		}
	}
}

func (inst *FFPoint) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		inst.setTopic()
	}
	if loopCount%retryCount == 0 {
		inst.checkStillExists()
	}
	val, null := inst.GetPayloadNull()
	var writeNull bool
	if null {
		writeNull = true
	} else {
		p, value, currentPri, err := parseCOV(val)
		if err == nil && p != nil {
			inst.lastPayload = p
			inst.WritePinFloat(node.Out, value, 2)
			inst.WritePinFloat(node.CurrentPriority, float64(currentPri))
			if inst.lastValue != value {
				inst.lastValue = value
				inst.lastUpdate = time.Now()
			} else {
				inst.WritePin(node.LastUpdated, ttime.TimeSince(inst.lastUpdate))
			}
		} else {
			writeNull = true
		}
	}
	if writeNull { // make sure we return some values
		inst.WritePinNull(node.Out)
		inst.WritePinNull(node.CurrentPriority)
		inst.WritePinNull(node.LastUpdated)
	}
}

func (inst *FFPoint) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
