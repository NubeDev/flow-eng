package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

type BV struct {
	*node.Spec
	objectID    points.ObjectID
	objectType  points.ObjectType
	pointUUID   string
	store       *points.Store
	application names.ApplicationName
}

const (
	object = "object"
)

func NewBV(body *node.Spec, opts *Bacnet) (node.Node, error) {
	var err error
	opts = bacnetOpts(opts)
	body, err = nodeDefault(body, bacnetBV, category, opts.Application)
	return &BV{
		body,
		0,
		points.BinaryVariable,
		"",
		opts.Store,
		opts.Application,
	}, err
}

func (inst *BV) setObjectId() {
	inst.objectID = points.ObjectID(inst.ReadPinAsInt(node.ObjectId))
}

func (inst *BV) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.setObjectId()
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		ioType := points.IoTypeNumber // TODO make a setting
		point := addPoint(inst.application, ioType, objectType, inst.objectID, isWriteable, isIO, true)
		point, err = inst.store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, points.AnalogInput, inst.objectID, inst.store)

}

func (inst *BV) Cleanup() {}
