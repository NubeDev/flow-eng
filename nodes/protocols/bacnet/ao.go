package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

type AO struct {
	*node.Spec
	objectID    points.ObjectID
	objectType  points.ObjectType
	pointUUID   string
	store       *points.Store
	application names.ApplicationName
}

func NewAO(body *node.Spec, opts *Bacnet) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, bacnetAO, category, opts.Application)
	return &AO{
		body,
		0,
		points.AnalogOutput,
		"",
		opts.Store,
		opts.Application,
	}, err
}
func (inst *AO) setObjectId() {
	inst.objectID = points.ObjectID(inst.ReadPinAsInt(node.ObjectId))
}
func (inst *AO) Process() {
	_, firstLoop := inst.Loop()
	if !firstLoop {
		inst.setObjectId()
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		ioType := points.IoTypeDigital // TODO make a setting
		point := addPoint(inst.application, ioType, objectType, inst.objectID, isWriteable, isIO, true)
		point, err = inst.store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, inst.objectID, inst.store)

}

func (inst *AO) Cleanup() {}
