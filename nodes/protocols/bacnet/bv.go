package bacnet

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

type BV struct {
	*node.Spec
	objectID   points.ObjectID
	objectType points.ObjectType
	pointUUID  string
}

const (
	object = "object"
)

func NewBV(body *node.Spec, store *points.Store) (node.Node, error) {
	var err error
	if store == nil {
		store = getStore()
	}
	body, err = nodeDefault(body, bacnetBV, category, store.GetApplication())
	return &BV{
		body,
		0,
		points.BinaryVariable,
		"",
	}, err
}

func (inst *BV) setObjectId() {
	inst.objectID = points.ObjectID(inst.ReadPinAsInt(node.ObjectId))
}

func (inst *BV) Process() {
	if !inst.OnStart {
		inst.setObjectId()
		store := getStore()
		objectType, isWriteable, _, err := getBacnetType(inst.Info.Name)
		ioType := points.IoTypeDigital
		point := addPoint(getApplication(), ioType, objectType, inst.objectID, isWriteable, false, true)
		point, err = store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, inst.objectID)
	fromFlow(inst, inst.objectID)
	inst.OnStart = true

}

func (inst *BV) Cleanup() {}
