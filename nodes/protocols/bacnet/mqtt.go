package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	rubixIO "github.com/NubeDev/flow-eng/helpers/rubixio"
	"github.com/NubeDev/flow-eng/nodes/protocols/points"
	log "github.com/sirupsen/logrus"
	"time"
)

func getTopic(msg interface{}) string {
	m := decode(msg)
	if m != nil {
		return m.Msg.Topic()
	}
	return ""
}

func decode(msg interface{}) *mqttbase.Message {
	m, ok := msg.(*mqttbase.Message)
	if ok {
		return m
	}
	return nil
}

//mqttRunner processMessage are the messages from the bacnet-server via the mqtt-broker
func (inst *Server) mqttSubRunner() {
	log.Info("start mqtt-sub-runner")
	go func() {
		msg, ok := inst.bus().Recv()
		if !ok {
			log.Errorf("bacnet-sub-runner malformed mqtt message")
		}
		decoded := decode(msg)
		var topic string
		if decoded != nil {
			topic = decoded.Msg.Topic()
		} else {
			log.Errorf("bacnet-sub-runner decoded malformed mqtt message")
		}
		if mqttbase.CheckBACnet(topic) {
			go inst.handleBacnet(msg)
		}
		if mqttbase.CheckRubixIO(topic) {
			go inst.handleRubixIO(decoded)
		}
	}()
}

func (inst *Server) handleBacnet(msg interface{}) {
	payload := points.NewPayload()
	err := payload.NewMessage(msg)
	if err != nil {
		log.Errorf("bacnet-sub-runner malformed mqtt message err:%s", err.Error())
		return
	}
	topic := payload.GetTopic()
	t, id := payload.GetObjectID()
	point := inst.db().GetPointByObject(t, id)
	if points.IsPri(topic) {
		value := payload.GetHighestPriority()
		log.Infof("mqtt-runner-subscribe point type:%s-%d value:%f", point.ObjectType, point.ObjectID, value.Value)
		inst.db().WritePointValue(point.UUID, value.Value)
	}
}

func (inst *Server) handleRubixIO(msg *mqttbase.Message) {
	r := &rubixIO.RubixIO{}
	inputs, err := r.DecodeInputs(msg.Msg.Payload())
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("RUBIX-INPUTS", inputs)

}

//mqttPubRunner send messages to the broker, as in read a modbus point and send it to the bacnet server
func (inst *Server) mqttPubRunner() {
	log.Info("start mqtt-pub-runner")
	for {
		for _, point := range inst.db().GetPoints() {
			if point.ToBacnetSyncPending {
				inst.mqttPublish(point)
				inst.db().UpdateBacnetSync(point.UUID, false)
			} else {
				//log.Infof("mqtt-runner-publish point skip as non cov type:%s-%d", point.ObjectType, point.ObjectID)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (inst *Server) mqttPublish(pnt *points.Point) {
	if pnt == nil {
		log.Errorf("bacnet-server-publish point can not be empty")
		return
	}
	objectType := pnt.ObjectType
	objectId := pnt.ObjectID
	value := pnt.ToBacnet
	obj, err := points.ObjectSwitcher(objectType)
	if err != nil {
		log.Error(err)
		return
	}
	topic := fmt.Sprintf("bacnet/%s/%d", obj, objectId)
	inst.client.Publish(value, topic)
}
