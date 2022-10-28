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
		s := inst.GetStore()
		children, ok := s.Get(inst.GetID())
		payloads := getPayloads(children, ok)
		for _, payload := range payloads {
			if payload.topic == fixTopic(message.Topic()) {
				n := inst.GetNode(payload.nodeUUID)
				if n != nil {
					//fmt.Println(string(message.Payload()), payload.nodeUUID, n.GetName())
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
					log.Errorf("flow-network-broker subscribe:%s err:%s", payload.topic, err.Error())
				} else {
					log.Infof("flow-network-broker subscribe:%s", payload.topic)
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
			log.Errorf("flow-network-broker subscribe:%s err:%s", topic, err.Error())
		} else {
			log.Infof("flow-network-broker subscribe:%s", topic)
		}
	}
}

func (inst *Network) pointsList() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		var points []*point
		err := json.Unmarshal(message.Payload(), &points)
		if err == nil {
			if points != nil {
				s := inst.GetStore()
				if s != nil {
					//fmt.Println(string(message.Payload()))
					s.Set(fmt.Sprintf("pointsList_%s", inst.GetID()), points, 0)
				}
			}
		} else {
			log.Errorf("failed to get flow-framework points list err:%s", err.Error())
		}
	}
	var topic = "rubix/platform/points/publish"
	if inst.mqttClient != nil {
		err := inst.mqttClient.Subscribe(topic, mqttQOS, callback)
		if err != nil {
			log.Errorf("flow-network-broker subscribe:%s err:%s", topic, err.Error())
		} else {
			log.Infof("flow-network-broker subscribe:%s", topic)
		}
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
		log.Errorf("flow-network-point publish faild to get point store")
		return
	}
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	for _, payload := range payloads {
		if !payload.isWriteable {
			continue
		}
		names := spitPointNames(payload.netDevPntNames)
		if len(names) != 4 {
			log.Errorf("flow-network-point publish failed to get point name")
			continue
		}
		n := inst.GetNode(payload.nodeUUID)
		if n == nil {
			log.Errorf("flow-network-point publish failed to get point node")
			continue
		}
		value, valueNull := n.ReadPinAsFloat(node.In16)
		if valueNull {
			log.Errorf("flow-network-point publish failed to get point input value")
			continue
		}
		priority := map[string]*float64{"_16": &value}
		pointWriter := &PointWriter{Priority: &priority}
		body := &MqttPoint{
			NetworkName: names[1],
			DeviceName:  names[2],
			PointName:   names[3],
			Priority:    pointWriter,
		}
		data, err := json.Marshal(body)
		if err != nil {
			log.Errorf("flow-network-broker publish err:%s", err.Error())
			continue
		}
		if inst.mqttClient != nil {
			updated, _ := n.InputUpdated(node.In16)
			if updated || loopCount == 2 {
				//fmt.Println("MQTT publish", string(data))
				err := inst.mqttClient.Publish(pointWriteTopic, mqttQOS, mqttRetain, data)
				if err != nil {
					log.Errorf("flow-network-broker publish err:%s", err.Error())
				}
			}
		}
	}
}
