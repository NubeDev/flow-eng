package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

func updateInputs(body node.Node, objectType points.ObjectType, id points.ObjectID) {
	v, _ := getStore().GetValueFromReadByObject(objectType, id)
	body.WritePin(node.Out, v)
	p := getStore().GetPointByObject(objectType, id)
	fmt.Println(getServer().client.IsConnected())
	getServer().mqttPublish(p)

}

func fromFlow(body node.Node) {
	objectType, isWriteable, _, err := getBacnetType(body.GetName())
	if err != nil {
		return
	}
	var ok bool
	var in14 *float64
	var in15 *float64
	if isWriteable {
		in14, ok = conversions.GetFloatPointer(body.ReadPin(node.In14))
		if !ok {
			log.Errorf("bacnet-server: failed to get node write value from node process")
		}
		in15, ok = conversions.GetFloatPointer(body.ReadPin(node.In15))
		if !ok {
			log.Errorf("bacnet-server: failed to get node write value from node process")
		}
		fmt.Println(in14, in15)
	}
	objectId, ok := conversions.GetInt(body.ReadPin(node.ObjectId))
	if !ok {
		log.Errorf("bacnet-server: failed to get node write value from node process")
	}
	if objectId == 0 {
		log.Errorf("bacnet-server: failed to get object-id from node process")
		return
	}
	createSync(nil, objectType, points.ObjectID(objectId), points.FromFlow, in14, in15)
}

func (inst *Server) fromBacnet(msg interface{}) {
	payload := points.NewPayload()
	err := payload.NewMessage(msg)
	if err != nil {
		log.Errorf("bacnet-sub-runner malformed mqtt message err:%s", err.Error())
		return
	}
	topic := payload.GetTopic()
	objectType, objectId := payload.GetObjectID()
	point := getStore().GetPointByObject(objectType, objectId)
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
		createSync(value, objectType, objectId, points.FromMqttPriory, nil, nil)
	}
}

func getToSync() points.SyncTo {
	app := getApplication()
	switch app {
	case applications.RubixIO:
		return points.ToRubixIO
	case applications.RubixIOAndModbus:
		return points.ToRubixIOModbus
	}
	return ""
}

// createSync can come from bacnet or the flow
func createSync(writeValue *points.PriArray, object points.ObjectType, id points.ObjectID, syncFrom points.SyncFrom, in14, in15 *float64) {
	point := getStore().GetPointByObject(object, id)
	if object == "" {
		log.Errorf("bacnet-server: object type type can not be empty")
	}
	if syncFrom == "" {
		log.Errorf("bacnet-server: get sync from can not be empty")
	}
	sync := getToSync()
	if sync == "" {
		log.Errorf("bacnet-server: get sync type can not be empty")
	}
	if point != nil {
		getStore().AddSync(point.UUID, writeValue, syncFrom, sync, getApplication())
		getStore().WritePointValue(point.UUID, writeValue, in14, in15)
	}

}
