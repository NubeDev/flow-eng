package flow

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (inst *Network) subscribeToEachPoint() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		// log.Infof("Flow Network subscribe(): %+v", string(message.Payload()))
		// log.Infof("Flow Network subscribe(): %+v", string(message.Topic()))
		s := inst.GetStore()
		children, ok := s.Get(inst.GetID())
		payloads := getPayloads(children, ok)
		for _, payload := range payloads {
			fixedTopic, pntUUID := fixTopic(message.Topic())
			payload.pointUUID = pntUUID
			if payload.topic == fixedTopic {
				n := inst.GetNode(payload.nodeUUID)
				if n != nil {
					n.SetPayload(&node.Payload{
						Any: message,
					})
				}
			}
		}
	}
	s := inst.GetStore()
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	inst.subscribeFailedPoints = true
	for _, payload := range payloads {
		if payload.topic != "" {
			if inst.mqttClient != nil {
				err := inst.mqttClient.Subscribe(payload.topic, mqttQOS, callback)
				if err != nil {
					log.Errorf("Flow Network subscribe(): %s err: %s", payload.topic, err.Error())
				} else {
					log.Infof("Flow Network subscribe(): %s", payload.topic)
					inst.subscribeFailedPoints = false
				}
			}
		}
	}
}

func (inst *Network) fetchPointsList() {
	var topic = "rubix/platform/points"
	inst.fetchPointResponseCount++
	if inst.mqttClient != nil {
		err := inst.mqttClient.Publish(topic, mqttQOS, false, "")
		if err != nil {
			log.Errorf("Flow Network fetchPointsList(): %s err: %s", topic, err.Error())
			inst.error = true
			inst.errorCode = errorFetchPointMQTTConnect
		} else {
			inst.error = false
			inst.errorCode = errorOk
		}
	} else {
		inst.error = true
		inst.errorCode = errorMQTTClientEmpty
	}
	if inst.fetchPointResponseCount > 5 {
		inst.error = true
		inst.errorCode = errorFailedFetchPoint
	}
}

// pointsList
// TODO: the points list topic gets a message on EVERY FF Point CRUD.  This produces MANY updates and many FF DB calls. Consider removing them.
func (inst *Network) pointsList() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		inst.fetchPointResponseCount = 0
		var points []*point
		err := json.Unmarshal(message.Payload(), &points)
		if err != nil {
			log.Errorf("failed to unmarshal flow-framework points list err: %s", err.Error())
			return
		}
		inst.pointsCount = len(points)
		log.Infof("Flow Network pointsList() points count: %d", inst.pointsCount)
		s := inst.GetStore()
		topic := fmt.Sprintf("pointsList_%s", inst.GetID())
		if s != nil {
			s.Set(topic, points, 0)
			inst.SetSubTitle(fmt.Sprintf("points count: %d", inst.pointsCount))
		} else {
			log.Errorf("failed to get flow-framework points store err: %s", err.Error())
		}
	}
	var topic = "rubix/platform/points/publish"
	inst.subscribeFailedPointsList = true
	if inst.mqttClient != nil {
		err := inst.mqttClient.Subscribe(topic, mqttQOS, callback)
		if err != nil {
			log.Errorf("Flow Network pointsList() :%s Subscribe() err: %s", topic, err.Error())
		} else {
			log.Infof("Flow Network pointsList() Subscribe() : %s", topic)
			inst.subscribeFailedPointsList = false
		}
	}
}

func (inst *Network) fetchSchedulesList() {
	var topic = "rubix/platform/schedules"
	if inst.mqttClient != nil {
		err := inst.mqttClient.Publish(topic, mqttQOS, false, "")
		if err != nil {
			log.Errorf("Flow Network fetchSchedulesList(): %s err: %s", topic, err.Error())
			inst.error = true
			inst.errorCode = errorFetchPointMQTTConnect
		} else {
			inst.error = false
			inst.errorCode = errorOk
		}
	} else {
		inst.error = true
		inst.errorCode = errorMQTTClientEmpty
	}
	if inst.fetchPointResponseCount > 5 {
		inst.error = true
		inst.errorCode = errorFailedFetchPoint
	}
}

// schedulesList
func (inst *Network) schedulesList() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		var schedules []*Schedule
		err := json.Unmarshal(message.Payload(), &schedules)
		if err != nil {
			log.Errorf("failed to unmarshal flow-framework schedules list err: %s", err.Error())
			return
		}
		log.Infof("Flow Network schedulesList() schedules count: %d", len(schedules))
		s := inst.GetStore()
		topic := fmt.Sprintf("schedulesList_%s", inst.GetID())
		if s != nil {
			s.Set(topic, schedules, 0)
		} else {
			log.Errorf("failed to get flow-framework schedules store err: %s", err.Error())
		}
	}
	var topic = "rubix/platform/schedules/publish"
	inst.subscribeFailedSchedulesList = true
	if inst.mqttClient != nil {
		err := inst.mqttClient.Subscribe(topic, mqttQOS, callback)
		if err != nil {
			log.Errorf("Flow Network schedulesList() :%s Subscribe() err: %s", topic, err.Error())

		} else {
			log.Infof("Flow Network schedulesList() Subscribe() : %s", topic)
			inst.subscribeFailedSchedulesList = false
		}
	}
}

const fetchSelectedPointsCOVTopic = "rubix/platform/points/cov/selected"

func (inst *Network) fetchAllPointValues() {
	// log.Infof("Flow Network fetchAllPointValues()")
	s := inst.GetStore()
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	var pointUUIDList []*MqttPoint
	for _, payload := range payloads {
		if payload.pointUUID != "" {
			pnt := MqttPoint{
				PointUUID: payload.pointUUID,
			}
			exists := false
			for _, val := range pointUUIDList {
				if val != nil && val.PointUUID == pnt.PointUUID {
					exists = true
					break
				}
			}
			if !exists {
				pointUUIDList = append(pointUUIDList, &pnt)
			}
		}
	}
	data, err := json.Marshal(pointUUIDList)
	if err != nil {
		log.Errorf("Flow Network fetchAllPointValues() json marshal err: %s", err.Error())
		return
	}
	if inst.mqttClient != nil {
		err := inst.mqttClient.Publish(fetchSelectedPointsCOVTopic, mqttQOS, mqttRetain, data)
		if err != nil {
			log.Errorf("Flow Network fetchAllPointValues() mqtt publish err: %s", err.Error())
		}
	} else {
		log.Errorf("Flow Network fetchAllPointValues() mqttClient not available")
	}
}

func spitPointNames(names string) []string {
	s := strings.Split(names, ":")
	if len(s) == 4 {
		return s
	}
	return nil
}

const pointWriteTopic = "rubix/platform/point/write"

func (inst *Network) publish(loopCount uint64) {
	s := inst.GetStore()
	if s == nil {
		log.Errorf("Flow Network publish() err: faild to get point store")
		return
	}
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	// log.Infof("Flow Network publish() point count to write: %d", len(payloads))

	for _, payload := range payloads {
		// fmt.Println(fmt.Sprintf("Flow Network publish() payload: %+v", payload))
		if !payload.isWriteable {
			continue
		}
		names := spitPointNames(payload.netDevPntNames)
		// log.Infof("spitPointNames(): %+v", names)
		if len(names) != 4 {
			log.Errorf("Flow Network publish() err: failed to get point name, incorrect lenght")
			continue
		}
		n := inst.GetNode(payload.nodeUUID)
		if n == nil {
			log.Errorf("Flow Network publish() err: failed to get point node")
			continue
		}

		republishLoop := false
		if loopCount%100 == 0 { // republish every 100 loops
			republishLoop = true
			// log.Infof("Flow Network publish(): REPUBLISH LOOP!")
		}

		priority := map[string]*float64{}
		// TODO: Replace this next bit with Priority Array evaluation
		if n.GetName() == flowPointWrite {
			ffPointWriteNode := n.(*FFPointWrite)
			priority = ffPointWriteNode.EvaluateInputsArray(republishLoop)
			if len(priority) <= 0 {
				// log.Infof("Flow Network publish() nothing to write name: %s", payload.netDevPntNames)
				continue
			}

			pointWriter := &PointWriter{Priority: &priority}
			body := &MqttPoint{
				NetworkName: names[1],
				DeviceName:  names[2],
				PointName:   names[3],
				Priority:    pointWriter,
			}
			data, err := json.Marshal(body)
			if err != nil {
				log.Errorf("Flow Network publish() err: %s", err.Error())
				continue
			}
			if inst.mqttClient != nil {
				err := inst.mqttClient.Publish(pointWriteTopic, mqttQOS, mqttRetain, data)
				if err != nil {
					log.Errorf("Flow Network publish() err: %s", err.Error())
				}
			} else {
				log.Errorf("Flow Network publish() mqtt client is empty")
			}
		}
	}
}
