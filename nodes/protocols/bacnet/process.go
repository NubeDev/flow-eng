package bacnetio

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
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
	var in14 *float64
	var in15 *float64
	if isWriteable {
		in14Value, in14Null := body.ReadPinAsFloat(node.In14)
		if in14Null {
			in14 = nil
		} else {
			in14 = float.New(in14Value)
		}
		in15Value, in15Null := body.ReadPinAsFloat(node.In15)
		if in15Null {
			in15 = nil
		} else {
			in15 = float.New(in15Value)
		}
	}
	if objectId == 0 {
		log.Errorf("bacnet-server: failed to get object-id from node process")
		return
	}
	store.CreateSync(nil, objectType, objectId, points.FromFlow, in14, in15)
}

func fromBacnet(msg interface{}, store *points.Store) error {
	payload := points.NewPayload()
	err := payload.NewMessage(msg)
	if err != nil {
		return err
	}
	topic := payload.GetTopic()
	objectType, objectId := payload.GetObjectID()
	point := store.GetPointByObject(objectType, objectId)
	if point == nil {
		return errors.New(fmt.Sprintf("mqtt-payload-priorty-array no point-found in store for type:%s-%d", objectType, objectId))
	}
	if topics.IsPri(topic) {
		value := payload.GetFullPriority()
		store.CreateSync(value, objectType, objectId, points.FromMqttPriory, nil, nil)
	}
	return nil
}
