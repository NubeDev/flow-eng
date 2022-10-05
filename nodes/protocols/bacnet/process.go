package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

/*
PROCESS
Is where a message will come from the flow or bacnet
These would only be messages that we need to write to, as in write output-point, to a modbus or edge-28 point
*/

// fromFlow is when a node has been written to from the wire sheet link, as in write a value @16
func fromFlow(body node.Node, objectId points.ObjectID, store *points.Store) {
	objectType, isWriteable, _, err := getBacnetType(body.GetName())
	if err != nil {
		return
	}
	var ok bool
	var in14 *float64
	var in15 *float64
	if isWriteable {
		in14, ok = conversions.GetFloatPointerOk(body.ReadPin(node.In14))
		if !ok {
			log.Errorf("bacnet-server: @14 failed to get node write value from node process")
		}
		in15, ok = conversions.GetFloatPointerOk(body.ReadPin(node.In15))
		if !ok {
			log.Errorf("bacnet-server:  @15 failed to get node write value from node process")
		}
	}
	if objectId == 0 {
		log.Errorf("bacnet-server: failed to get object-id from node process")
		return
	}
	store.CreateSync(nil, objectType, objectId, points.FromFlow, in14, in15)
}

func fromBacnet(msg interface{}, store *points.Store) {
	payload := points.NewPayload()
	err := payload.NewMessage(msg)
	if err != nil {
		log.Errorf("bacnet-sub-runner malformed mqtt message err:%s", err.Error())
		return
	}
	topic := payload.GetTopic()
	objectType, objectId := payload.GetObjectID()
	point := store.GetPointByObject(objectType, objectId)
	if point == nil {
		log.Errorf("mqtt-payload-priorty-array no point-found in store for type:%s-%d", objectType, objectId)
		return
	}
	if topics.IsPri(topic) {
		value := payload.GetFullPriority()
		highest := payload.GetHighestPriority()
		if highest != nil {
			log.Infof("mqtt-runner-subscribe point type:%s-%d value:%f", point.ObjectType, point.ObjectID, highest.Value)
		}
		store.CreateSync(value, objectType, objectId, points.FromMqttPriory, nil, nil)
	}
}
