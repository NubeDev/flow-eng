package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
)

type BO struct {
	*node.Spec
	objectID      points.ObjectID
	objectType    points.ObjectType
	pointUUID     string
	store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

func NewBO(body *node.Spec, opts *Bacnet) (node.Node, error) {
	var err error
	opts = bacnetOpts(opts)
	body, err = nodeDefault(body, bacnetBO, category, opts.Application)
	flowOptions := &toFlowOptions{}
	return &BO{
		body,
		0,
		points.BinaryOutput,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
		flowOptions,
	}, err
}

func (inst *BO) setName() {
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

func (inst *BO) setObjectId() {
	id, _ := inst.ReadPinAsInt(node.ObjectId)
	inst.objectID = points.ObjectID(id)
}

func (inst *BO) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.setObjectId()
		inst.setName()
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		ioType := points.IoTypeDigital
		point := addPoint(ioType, objectType, inst.objectID, isWriteable, isIO, true, inst.application)
		point, err = inst.store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, points.BinaryOutput, inst.objectID, inst.store, inst.toFlowOptions)
	fromFlow(inst, inst.objectID, inst.store)
}
