package ros

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
)

type ROSPoint struct {
	*node.Spec
	topic           string
	lastPayload     *covPayload
	lastValue       *float64
	lastUpdate      time.Time
	networkNodeUUID string
}

func NewROSPoint(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, rosPoint, Category)
	inputs := node.BuildInputs()
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	currentPriority := node.BuildOutput(node.CurrentPriority, node.TypeFloat, nil, body.Outputs)
	lastUpdated := node.BuildOutput(node.LastUpdated, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out, currentPriority, lastUpdated)
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	pnt := &ROSPoint{body, "", nil, nil, time.Now(), ""}
	return pnt, nil
}

// only one rubix-network can be added
// we need to get the node uuid of the rubix-network node for each point node

func (inst *ROSPoint) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		uuid, err := inst.getRubixNetworkUUID()
		if err != nil {
			log.Error(err)
		}
		inst.networkNodeUUID = uuid
		inst.setTopic()
	}
	if loopCount%retryCount == 0 {
		inst.setTopic()
	}

	if inst.checkStillExists() {
		val, lastUpdated, null := inst.GetPayloadNull()
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
				inst.lastUpdate = lastUpdated

				if currentPri == nil {
					inst.WritePinNull(node.CurrentPriority)
				} else {
					inst.WritePinFloat(node.CurrentPriority, float64(*currentPri))
				}

				inst.WritePin(node.LastUpdated, ttime.TimePretty(inst.lastUpdate))
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

// set add this point to the store
func (inst *ROSPoint) set() {
	// log.Infof("FLOW POINT: set() topic: %+v", inst.topic)
	// log.Infof("FLOW POINT: set() STORE: %+v", inst.GetStore().All())
	// log.Infof("FLOW POINT: set() STORE Object: %+v", inst.GetStore().All()[inst.GetParentId()].Object)
	s := inst.GetStore()
	if s == nil {
		return
	}
	parentId := inst.networkNodeUUID
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

func (inst *ROSPoint) checkStillExists() bool {
	nodeStatus := inst.GetStatus()
	if nodeStatus.WaringMessage == pointError { // This is set by the subscribeToMissingPoints() runner
		return false
	} else {
		return true
	}
}

func (inst *ROSPoint) setTopic() {
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

func (inst *ROSPoint) getRubixNetworkUUID() (string, error) {
	nodes := inst.GetNodesByType(rosNetwork)
	if len(nodes) == 0 {
		return "", errors.New("no rubix-network node has been added")
	}
	if len(nodes) > 1 {
		return "", errors.New("only one rubix-network node can be been added, please add one")
	}
	return nodes[0].GetID(), nil
}

func (inst *ROSPoint) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
