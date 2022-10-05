package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
	"time"
)

// mqttPubRunner send messages to the broker, as in read a modbus point and send it to the bacnet server
func (inst *Server) writeRunner() {
	log.Info("start mqtt-pub-runner")
	for {
		for _, point := range inst.store.GetPoints() {
			inst.mqttPublishPV(point)
		}
		time.Sleep(runnerDelay * time.Millisecond)
	}
}

// mqttPublish example for future MQTT write to the bacnet-server
func (inst *Server) mqttPublishPV(pnt *points.Point) {
	if pnt == nil {
		log.Errorf("bacnet-server-publish point can not be empty")
		return
	}
	objectType := pnt.ObjectType
	objectId := pnt.ObjectID
	value := pnt.WriteValue
	obj, err := points.ObjectSwitcher(objectType)
	if err != nil {
		log.Error(err)
		return
	}
	v := points.GetHighest(value)
	topic := fmt.Sprintf("bacnet/%s/%d/write/pv", obj, objectId) // bacnet/ao/1/write/pv
	if v != nil {
		err = inst.clients.mqttClient.Publish(topic, mqttclient.AtMostOnce, true, fmt.Sprintf("%f", v.Value))
		if err != nil {
			log.Errorf("bacnet-server: mqtt publish err: %s", err.Error())
			return
		}
	}

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

*/
