package bacnet

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

type AO struct {
	*node.Spec
	objectID   points.ObjectID
	objectType points.ObjectType
	pointUUID  string
}

func NewAO(body *node.Spec, store *points.Store) (node.Node, error) {
	var err error
	if store == nil {
		store = getStore()
	}
	body, err = nodeDefault(body, bacnetAO, category, store.GetApplication())
	return &AO{
		body,
		0,
		points.AnalogOutput,
		"",
	}, err
}
func (inst *AO) setObjectId() {
	inst.objectID = points.ObjectID(inst.ReadPinAsInt(node.ObjectId))
}
func (inst *AO) Process() {
	if !inst.OnStart {
		inst.setObjectId()
		store := getStore()
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		ioType := points.IoTypeDigital // TODO make a setting
		point := addPoint(getApplication(), ioType, objectType, inst.objectID, isWriteable, isIO, true)
		point, err = store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
	}
	toFlow(inst, inst.objectID)
	fromFlow(inst, inst.objectID)
	inst.OnStart = true

}

func (inst *AO) Cleanup() {}
