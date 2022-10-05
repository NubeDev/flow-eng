package bacnet

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
	objectID    points.ObjectID
	objectType  points.ObjectType
	pointUUID   string
	store       *points.Store
	application names.ApplicationName
	mqttClient  *mqttclient.Client
}

func NewAI(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	var err error
	body, err = nodeDefault(body, bacnetAI, category, opts.Application)
	return &AI{
		body,
		0,
		points.AnalogInput,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
	}, err
}

func (inst *AI) setName() {
	// bacnet/ao/1/write/name
	name := inst.ReadPinAsString(node.Name)
	if name != "" {
		topic := fmt.Sprintf("%s/write/name", topicBuilder(inst.objectType, inst.objectID))
		err := inst.mqttClient.Publish(topic, mqttclient.AtMostOnce, true, name)
		if err != nil {
			return
		}
	}

}

func (inst *AI) setObjectId() {
	inst.objectID = points.ObjectID(inst.ReadPinAsInt(node.ObjectId))
}

func (inst *AI) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.setObjectId()
		inst.setName()
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		ioType := points.IoTypeTemp // TODO make a setting
		point := addPoint(inst.application, ioType, objectType, inst.objectID, isWriteable, isIO, true)
		point, err = inst.store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, points.AnalogInput, inst.objectID, inst.store)
}

func (inst *AI) Cleanup() {}
