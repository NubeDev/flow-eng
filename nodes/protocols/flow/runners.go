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
				n.SetPayload(&node.Payload{
					Any: message,
					//String: str.New(string(message.Payload())),
				})
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

func (inst *Network) pointsList() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		var points []*point
		err := json.Unmarshal(message.Payload(), &points)
		if err == nil {
			if points != nil {
				s := inst.GetStore()
				s.Set(fmt.Sprintf("pointsList_%s", inst.GetID()), points, 0)
			}
		} else {
			log.Errorf("failed to get flow-framework points list err:%s", err.Error())
		}
	}
	var topic = "+/+/+/+/+/+/rubix/points/value/points"
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

const fetchPointsTopicWrite = "rubix/platform/points/write"

func (inst *Network) publish(loopCount uint64) {
	s := inst.GetStore()
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	for _, payload := range payloads {
		if !payload.isWriteable {
			continue
		}
		names := spitPointNames(payload.netDevPntNames)
		if len(names) != 4 {
			continue
		}
		n := inst.GetNode(payload.nodeUUID)
		value, valueNull := n.ReadPinAsFloat(node.In16)
		if valueNull {
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
				err := inst.mqttClient.Publish(fetchPointsTopicWrite, mqttQOS, mqttRetain, string(data))
				if err != nil {
					log.Errorf("flow-network-broker publish err:%s", err.Error())
				}
			}
		}
	}
}
