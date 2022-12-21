package flow

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (inst *Network) subscribe() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		// log.Infof("Flow Network subscribe(): %+v", message)
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
	for _, payload := range payloads {
		if payload.topic != "" {
			if inst.mqttClient != nil {
				err := inst.mqttClient.Subscribe(payload.topic, mqttQOS, callback)
				if err != nil {
					log.Errorf("Flow Network subscribe(): %s err: %s", payload.topic, err.Error())
				} else {
					log.Infof("Flow Network subscribe(): %s", payload.topic)
				}
			}
		}
	}
}

func (inst *Network) fetchPointsList() {
	var topic = "rubix/platform/points"
	if inst.mqttClient != nil {
		err := inst.mqttClient.Publish(topic, mqttQOS, false, "")
		if err != nil {
			log.Errorf("Flow Network fetchPointsList(): %s err: %s", topic, err.Error())
		}
	}
}

// TODO: the points list topic gets a message on EVERY FF Point CRUD.  This produces MANY updates and many FF DB calls. Consider removing them.
func (inst *Network) pointsList() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		var points []*point
		err := json.Unmarshal(message.Payload(), &points)
		// log.Infof("Flow Network pointsList() points: %+v", points)
		if err == nil {
			if points != nil {
				s := inst.GetStore()
				if s != nil {
					// fmt.Println(string(message.Payload()))
					s.Set(fmt.Sprintf("pointsList_%s", inst.GetID()), points, 0)
				}
			}
		} else {
			log.Errorf("failed to get flow-framework points list err: %s", err.Error())
		}
	}
	var topic = "rubix/platform/points/publish"
	if inst.mqttClient != nil {
		err := inst.mqttClient.Subscribe(topic, mqttQOS, callback)
		if err != nil {
			log.Errorf("Flow Network pointsList() :%s Subscribe() err: %s", topic, err.Error())
		} else {
			log.Infof("Flow Network pointsList() Subscribe() : %s", topic)
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
	for _, payload := range payloads {
		// fmt.Println(fmt.Sprintf("Flow Network publish() payload: %+v", payload))
		if !payload.isWriteable {
			continue
		}
		names := spitPointNames(payload.netDevPntNames)
		// log.Infof("spitPointNames(): %+v", names)
		if len(names) != 4 {
			log.Errorf("Flow Network publish() err: failed to get point name")
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
				// fmt.Println("MQTT publish", string(data))
				err := inst.mqttClient.Publish(pointWriteTopic, mqttQOS, mqttRetain, data)
				if err != nil {
					log.Errorf("Flow Network publish() err: %s", err.Error())
				}
			}
		}
	}
}
