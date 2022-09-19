package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
	log "github.com/sirupsen/logrus"
	"time"
)

//mqttRunner processMessage are the messages from the bacnet-server via the mqtt-broker
func (inst *Server) mqttSubRunner() {
	log.Info("start mqtt-sub-runner")
	go func() {
		msg, ok := inst.bus().Recv()
		if ok {
			payload := bstore.NewPayload()
			err := payload.NewMessage(msg)
			if err != nil {
				log.Errorf("bacnet-sub-runner malformed mqtt message err:%s", err.Error())
				return
			}
			t, id := payload.GetObject()
			pnt := inst.db().GetPointByObject(t, id)
			inst.db().WritePointValue(pnt.UUID, float.NonNil(payload.GetPresentValue()))
		}
	}()
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
				log.Infof("mqtt-runner-publish point skip as non cov type:%s-%d", point.ObjectType, point.ObjectID)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (inst *Server) mqttPublish(pnt *bstore.Point) {
	if pnt == nil {
		log.Errorf("bacnet-server-publish point can not be empty")
		return
	}
	objectType := pnt.ObjectType
	objectId := pnt.ObjectID
	value := pnt.ToBacnet
	topic := fmt.Sprintf("bacnet/%s/%d", objectType, objectId)
	inst.client.Publish(value, topic)
}
