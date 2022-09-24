package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

type AI struct {
	*node.Spec
	onStart    bool
	objectID   points.ObjectID
	objectType points.ObjectType
	pointUUID  string
}

func NewAI(body *node.Spec, store *points.Store) (node.Node, error) {
	var err error
	if store == nil {
		store = getStore()
	}
	body, err = nodeDefault(body, bacnetAI, category, store.GetApplication())

	return &AI{
		body,
		false,
		0,
		points.AnalogInput,
		"",
	}, err
}

func (inst *AI) setObjectId() {
	id, ok := conversions.GetInt(inst.ReadPin(node.ObjectId))
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}

func (inst *AI) Process() {
	if !inst.onStart {
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
	inst.onStart = true
}

func (inst *AI) Cleanup() {}
