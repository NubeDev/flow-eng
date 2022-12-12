package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
)

type AO struct {
	*node.Spec
	objectID      points.ObjectID
	objectType    points.ObjectType
	pointUUID     string
	store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

func NewAO(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	var err error
	body, err = nodeDefault(body, bacnetAO, category, opts.Application)
	body.SetSchema(buildSchemaUO())
	flowOptions := &toFlowOptions{}
	return &AO{
		body,
		0,
		points.AnalogOutput,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
		flowOptions,
	}, err
}

func (inst *AO) setName() {
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

func (inst *AO) setObjectId() {
	id, _ := inst.ReadPinAsInt(node.ObjectId)
	inst.objectID = points.ObjectID(id)
}
func (inst *AO) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.setObjectId()
		go inst.setName()
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		settings, err := getSettings(inst.GetSettings())
		ioType := settings.Io
		if ioType == "" {
			ioType = string(points.IoTypeVolts)
		}
		inst.toFlowOptions.precision = settings.Decimal
		point := addPoint(points.IoType(ioType), objectType, inst.objectID, isWriteable, isIO, true, inst.application)
		point, err = inst.store.AddPoint(point, false)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, points.AnalogOutput, inst.objectID, inst.store, inst.toFlowOptions)
	fromFlow(inst, inst.objectID, inst.store)
}
