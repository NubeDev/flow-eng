package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
)

type AI struct {
	*node.Spec
	objectID   points.ObjectID
	objectType points.ObjectType
	pointUUID  string
	//store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

var store *points.Store

func NewAI(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	var err error
	body, err = nodeDefault(body, bacnetAI, category, opts.Application)
	body.SetSchema(buildSchemaUI())
	flowOptions := &toFlowOptions{}
	store = opts.Store
	return &AI{
		body,
		0,
		points.AnalogInput,
		"",
		//opts.Store,
		opts.Application,
		opts.MqttClient,
		flowOptions,
	}, err
}

func (inst *AI) setName() {
	name, null := inst.ReadPinAsString(node.Name)
	if null {
		name = fmt.Sprintf("%s_%d", inst.objectType, inst.objectID)
	}
	topic := fmt.Sprintf("%s/write/name", topicBuilder(inst.objectType, inst.objectID))
	payload := buildPayload(name, 0)
	if payload != "" {
		err := inst.mqttClient.Publish(topic, mqttQOS, mqttRetain, payload)
		if err != nil {
		}
	}
}

func (inst *AI) setObjectId() {
	id, _ := inst.ReadPinAsInt(node.ObjectId)
	inst.objectID = points.ObjectID(id)
}

func (inst *AI) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.setObjectId()
		inst.setName()
		settings, err := getSettings(inst.GetSettings())
		ioType := settings.Io
		if ioType == "" {
			ioType = string(points.IoTypeVolts)
		}
		inst.toFlowOptions.precision = settings.Decimal
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		point := addPoint(points.IoType(ioType), objectType, inst.objectID, isWriteable, isIO, true, inst.application)
		fmt.Println(22222222)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d err:%s", objectType, inst.objectID, err.Error())
			return
		}
		fmt.Println(111111111)
		s := inst.GetStore()
		if s == nil {
			log.Errorf("bacnet-server add new point failed to get store type:%s-%d err:%s", objectType, inst.objectID, err.Error())
			return
		}
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!11", setUUID(points.AnalogInput, inst.objectID))
		s.Set(setUUID(points.AnalogInput, inst.objectID), point, 0)
		//point, err = store.AddPoint(point, true)
	}
	toFlow(inst, points.AnalogInput, inst.objectID, store, inst.toFlowOptions)

	s := inst.GetStore()
	if s == nil {
		return
	}
	//parentId := inst.GetParentId()
	//nodeUUID := inst.GetID()

	//s.Set("123", 123, 0)
	//
	//d, ok := s.Get("123")
	//
	////fmt.Println(parentId)
	//fmt.Println(nodeUUID)
	//fmt.Println(ok)
	//fmt.Println(d)

}
