package bacnetio

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (inst *Server) mqttReconnect() {
	var err error
	inst.pingLock = true
	err = inst.clients.mqttClient.Ping()
	if err != nil {
		log.Errorf("bacnet-server mqtt ping failed")
		inst.pingFailed = true
		err = inst.clients.mqttClient.Connect()
		if err != nil {
			log.Errorf("bacnet-server failed to reconnect with mqtt broker")
			inst.reconnectedOk = false
			inst.setWaringMessage("bacnet-server failed to reconnect with mqtt broker", false)
		} else {
			inst.reconnectedOk = true
			inst.setWaringMessage("", true)
		}
	}
	inst.pingLock = false
}

func (inst *Server) subscribeToBacnetServer() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		mes := &topics.Message{UUID: helpers.ShortUUID("bus"), Msg: message}
		if topics.IsPri(message.Topic()) {
			log.Infof("mqtt-bacnet message from server topic: %s -> value: %s", mes.Msg.Topic(), string(mes.Msg.Payload()))
			err := inst.fromBacnet(mes)
			if err != nil {
				if err.Error() != points.ErrStopMQTTLoop {
					log.Error(err)
				}
			}
		}
	}
	objsOuts := []string{"ao", "av", "bv"}
	for _, obj := range objsOuts {
		topic := fmt.Sprintf("%s/+/pri", topicObjectBuilder(points.ObjectType(obj)))
		err := inst.clients.mqttClient.Subscribe(topic, mqttQOS, callback)
		if err != nil {
			log.Errorf("bacnet-server mqtt:%s", err.Error())
		}
	}
	inst.pingFailed = false

}

func (inst *Server) subscribe() {
	inst.subscribeToBacnetServer()
}

// mqttPublish example for future MQTT write to the bacnet-server
func (inst *Server) mqttPublishPV(point *points.Point) error {
	if point == nil {
		return errors.New("bacnet-server-publish point can not be empty")
	}
	objectType := point.ObjectType
	objectId := point.ObjectID
	obj, err := points.ObjectSwitcher(objectType)
	if err != nil {
		log.Error(err)
		return err
	}

	updatePoint := false

	if point.IsWriteable {
		if point.PendingMQTTPublish {
			payload := ""
			topic := ""

			if false { // THIS WOULD BE USED FOR OUTPUT POINTS THAT DON'T USE THE PRIORITY ARRAY (ONLY PRESENT VALUE)

			} else {
				in14, in15 := points.GetWriteArrayValues(point.WriteValue)

				if in14 == nil {
					topic = fmt.Sprintf("bacnet/%s/%d/write/pri/14", obj, objectId)
					payload = buildPayload("null", 0)
				} else {
					topic = fmt.Sprintf("bacnet/%s/%d/write/pri/14", obj, objectId)
					payload = buildPayload("", *point.WriteValue.P14)
				}

				log.Infof("mqtt-bacnet publish (bacnet variable) topic: %s -> value: %s", topic, payload)
				if payload != "" {
					err = inst.clients.mqttClient.Publish(topic, mqttQOS, mqttRetain, payload) // write @14
					if err != nil {
						log.Errorf("bacnet-server: mqtt publish (bacnet variable) err: %s", err.Error())
						return err
					} else {
						point.PendingMQTTPublish = false
						updatePoint = true
					}
				}

				if in15 == nil {
					topic = fmt.Sprintf("bacnet/%s/%d/write/pri/15", obj, objectId)
					payload = buildPayload("null", 0)
				} else {
					topic = fmt.Sprintf("bacnet/%s/%d/write/pri/15", obj, objectId)
					payload = buildPayload("", *point.WriteValue.P15)
				}
				log.Infof("mqtt-bacnet publish (bacnet variable) topic: %s -> value: %s", topic, payload)
				if payload != "" {
					err = inst.clients.mqttClient.Publish(topic, mqttQOS, mqttRetain, payload) // write @15
					if err != nil {
						log.Errorf("bacnet-server: mqtt publish (bacnet variable) err: %s", err.Error())
						return err
					} else {
						point.PendingMQTTPublish = false
						updatePoint = true
					}
				}
			}
		}
	} else {
		if point.PendingMQTTPublish {
			payload := buildPayload("null", 0)
			if point.PresentValue != nil {
				payload = buildPayload("", *point.PresentValue)
			}
			topic := fmt.Sprintf("bacnet/%s/%d/write/pv", obj, objectId) // bacnet/ai/1/write/pv
			log.Infof("mqtt-bacnet publish (bacnet input) topic: %s -> value: %s", topic, payload)
			if payload != "" {
				err = inst.clients.mqttClient.Publish(topic, mqttQOS, mqttRetain, payload)
				if err != nil {
					log.Errorf("bacnet-server: mqtt publish (bacnet input) err: %s", err.Error())
					return err
				} else {
					point.PendingMQTTPublish = false
					updatePoint = true
				}
			}
		}
	}
	if updatePoint {
		inst.updatePoint(objectType, objectId, point)
	}
	return nil
}

// mqttPublishNames write the point names to the server
func (inst *Server) mqttPublishNames(point *points.Point) {
	if point == nil {
		log.Error("failed to find point to publish bacnet name to server")
		return
	}
	objectType := point.ObjectType
	objectId := point.ObjectID
	obj, err := points.ObjectSwitcher(objectType)
	if err != nil {
		log.Error(err)
		return
	}
	topic := fmt.Sprintf("bacnet/%s/%d/write/name", obj, objectId) // bacnet/ao/1/write/pv
	// topic := fmt.Sprintf("bacnet/%s/%d/write/pri/15", obj, objectId) // bacnet/ao/1/write/pv
	name := point.Name
	if name == "" {
		name = fmt.Sprintf("%s_%d", obj, objectId)
		name = strings.ToUpper(name)
	}
	payload := buildPayloadName(name)
	log.Infof("mqtt-bacnet publish name topic: %s -> value: %s", topic, payload)
	if payload != "" {
		err = inst.clients.mqttClient.Publish(topic, mqttQOS, mqttRetain, payload)
		if err != nil {
			log.Errorf("bacnet-server: mqtt publish err: %s", err.Error())
			return
		}
	}
	return
}

func getTopic(msg interface{}) string {
	m := decode(msg)
	if m != nil {
		return m.Msg.Topic()
	}
	return ""
}

func decode(msg interface{}) *topics.Message {
	m, ok := msg.(*topics.Message)
	if ok {
		return m
	}
	return nil
}

func (inst *Server) fromBacnet(msg interface{}) error {
	var err error
	payload := points.NewPayload()
	if payload == nil {
		return errors.New("mqtt-payload message from bacnet-server failed to get payload")
	}
	if !inst.firstMessageFromBacnet { // this is to try and get the bacnet-server's last value sent from another bacnet master
		err = payload.NewMessage(msg, true)
		if err == nil {
			// inst.firstMessageFromBacnet = true // Aidan: disabled as this logic dosnt work, reason is that flo-eng wasn't using the priority array from driver-bacnet
		}
	} else {
		err = payload.NewMessage(msg, false)
	}
	if err != nil {
		return err
	}
	topic := payload.GetTopic()
	objectType, objectId := payload.GetObjectID()
	point, _ := inst.getPoint(objectType, objectId)
	if point == nil {
		return errors.New(fmt.Sprintf("mqtt-payload-priorty-array no point-found in store for type:%s-%d", objectType, objectId))
	}
	if topics.IsPri(topic) {
		value := payload.GetFullPriority()
		err := inst.updateFromBACnet(objectType, objectId, value)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
update name

topic
bacnet/ao/0/write/name

{"value" : "ao_name0", "uuid" : "123456"}

to update present value of AO

topic
bacnet/ao/0/write/pv

json payload
{"value" : "100.50", "uuid" : "123456"}


Priority Array topic formats:
```
bacnet/object/address/write/pri/priority_index
bacnet/object/address/write/pri/priority_index/all
```

Write 50.20 into analog object (ao) at instance (address) 1 at priority index 10
```
topic: bacnet/ao/1/write/pri/10
json payload: {"value" : "10.50", "uuid" : "123456"}
```

Write 99.99 into analog object (ao) at instance (address) 1 to all priority slots
```
topic: bacnet/ao/1/write/pri/16/all
json payload: {"value" : "99.99", "uuid" : "123456"}
```

Reset all priority slots of analog object (ao) at instance (address) 1
```
topic: bacnet/ao/1/write/pri/16/all
json payload: {"value" : "null", "uuid" : "123456"}

*/
