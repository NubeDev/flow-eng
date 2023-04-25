package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
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
	lastValue   *float64
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
	pnt := &FFPoint{body, "", nil, nil, time.Now()}
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

func (inst *FFPoint) checkStillExists() bool {
	nodeStatus := inst.GetStatus()
	if nodeStatus.WaringMessage == pointError { // This is set by the subscribeToMissingPoints() runner
		return false
	} else {
		return true
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
		inst.setTopic()
	}

	if inst.checkStillExists() {
		val, null := inst.GetPayloadNull()
		// fmt.Println(fmt.Sprintf("FLOW POINT Process() payload: val: %+v,  null: %v", val, null))
		var writeNull bool
		if null {
			writeNull = true
		} else {
			p, value, currentPri, err := parseCOV(val)
			if err == nil && p != nil {
				inst.lastPayload = p
				if value == nil {
					inst.WritePinNull(node.Out)
					if inst.lastValue != nil {
						inst.lastValue = nil
					}
				} else {
					inst.WritePinFloat(node.Out, *value, 2)
					/* THIS IS USED TO SHOW LAST UPDATE AS LAST VALUE CHANGE
					if inst.lastValue == nil || *inst.lastValue != *value {
						inst.lastValue = float.New(*value)
						inst.lastUpdate = time.Now()
					}
					*/
				}
				// THIS IS USED TO SHOW THE LAST UPDATED AS THE LAST VALUE FETCH
				inst.lastValue = float.New(*value)
				inst.lastUpdate = time.Now()

				if currentPri == nil {
					inst.WritePinNull(node.CurrentPriority)
				} else {
					inst.WritePinFloat(node.CurrentPriority, float64(*currentPri))
				}

				inst.WritePin(node.LastUpdated, ttime.TimeSince(inst.lastUpdate))
			} else {
				if err != nil {
					fmt.Println(fmt.Sprintf("FLOW POINT Process() parseCOV() err: %v", err))
				}
				writeNull = true
			}
		}
		if writeNull { // make sure we return some values
			inst.WritePinNull(node.Out)
			inst.WritePinNull(node.CurrentPriority)
			inst.WritePinNull(node.LastUpdated)
		}
	} else {
		inst.WritePinNull(node.Out)
		inst.WritePinNull(node.CurrentPriority)
		inst.WritePinNull(node.LastUpdated)
	}

}

func (inst *FFPoint) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
