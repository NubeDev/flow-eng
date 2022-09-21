package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

func (inst *Server) handleBacnet(msg interface{}) {
	payload := points.NewPayload()
	err := payload.NewMessage(msg)
	if err != nil {
		log.Errorf("bacnet-sub-runner malformed mqtt message err:%s", err.Error())
		return
	}
	topic := payload.GetTopic()
	t, id := payload.GetObjectID()
	point := getStore().GetPointByObject(t, id)
	if point == nil {
		log.Errorf("mqtt-payload-priorty-array no point-found in store for type:%s-%d", t, id)
		return
	}
	if topics.IsPri(topic) {
		value := payload.GetHighestPriority()
		log.Infof("mqtt-runner-subscribe point type:%s-%d value:%f", point.ObjectType, point.ObjectID, value.Value)
		getStore().WritePointValue(point.UUID, value.Value)
	}

}
