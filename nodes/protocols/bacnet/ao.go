package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

type AO struct {
	*node.Spec
	onStart    bool
	objectID   points.ObjectID
	objectType points.ObjectType
	pointUUID  string
}

func NewAO(body *node.Spec) (node.Node, error) {
	var err error
	store := getStore()
	body, err, point := nodeDefault(body, bacnetAO, category, store.GetApplication())
	var pointUUID string
	if point != nil {
		pointUUID = point.UUID
	}
	return &AO{
		body,
		false,
		0,
		points.AnalogOutput,
		pointUUID,
	}, err
}
func (inst *AO) setObjectId() {
	id, ok := conversions.GetInt(inst.ReadPin(node.ObjectId))
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}
func (inst *AO) Process() {
	if !inst.onStart {
		inst.setObjectId()
	}
	toFlow(inst, inst.objectID)
	fromFlow(inst, inst.objectID)
	inst.onStart = true
}

func (inst *AO) Cleanup() {}
