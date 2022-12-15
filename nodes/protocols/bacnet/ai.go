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
	objectID      points.ObjectID
	objectType    points.ObjectType
	pointUUID     string
	store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

//var store *points.Store

func NewAI(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	var err error
	body, err = nodeDefault(body, bacnetAI, category, opts.Application)
	body.SetSchema(buildSchemaUI())
	flowOptions := &toFlowOptions{}
	//store = opts.Store
	return &AI{
		body,
		0,
		points.AnalogInput,
		"",
		opts.Store,
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
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d err:%s", objectType, inst.objectID, err.Error())
			return
		}
		s := inst.GetStore()
		if s == nil {
			log.Errorf("bacnet-server add new point failed to get store type:%s-%d err:%s", objectType, inst.objectID, err.Error())
			return
		}
		point, err = inst.store.AddPoint(point, true)
		s.Set(setUUID(inst.GetParentId(), points.AnalogInput, inst.objectID), point, 0)
	}
	s := inst.GetStore()
	if s == nil {
		return
	}

	pv, err := inst.getPV(points.AnalogInput, inst.objectID)
	if err != nil {
		return
	}
	inst.WritePinFloat(node.Out, pv, 2)

}

func (inst *AI) getPV(objType points.ObjectType, id points.ObjectID) (float64, error) {
	pnt, ok := inst.getPoint(objType, id)
	if ok {
		return pnt.PresentValue, nil
	}
	return 0, nil
}

func (inst *AI) getPoint(objType points.ObjectType, id points.ObjectID) (*points.Point, bool) {
	s := inst.GetStore()
	if s == nil {
		return nil, false
	}
	d, ok := s.Get(setUUID(inst.ParentId, objType, id))
	if ok {
		return d.(*points.Point), true
	}
	return nil, false
}
