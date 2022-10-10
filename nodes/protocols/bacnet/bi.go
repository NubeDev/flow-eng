package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
)

type BI struct {
	*node.Spec
	objectID    points.ObjectID
	objectType  points.ObjectType
	pointUUID   string
	store       *points.Store
	application names.ApplicationName
	mqttClient  *mqttclient.Client
}

func NewBI(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	var err error
	body, err = nodeDefault(body, bacnetBI, category, opts.Application)
	body.SetSchema(buildSchemaUI())
	return &BI{
		body,
		0,
		points.BinaryInput,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
	}, err
}

func (inst *BI) setName() {
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

func (inst *BI) setObjectId() {
	id, _ := inst.ReadPinAsInt(node.ObjectId)
	inst.objectID = points.ObjectID(id)
}

func (inst *BI) Process() {
	_, firstLoop := inst.Loop()

	if firstLoop {
		inst.setObjectId()
		inst.setName()
		ioType := points.IoTypeDigital
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		point := addPoint(ioType, objectType, inst.objectID, isWriteable, isIO, true)
		point, err = inst.store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, points.BinaryInput, inst.objectID, inst.store)
}
