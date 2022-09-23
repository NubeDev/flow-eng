package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

type AI struct {
	*node.Spec
	onStart    bool
	objectID   points.ObjectID
	objectType points.ObjectType
	pointUUID  string
}

func NewAI(body *node.Spec) (node.Node, error) {
	var err error
	store := getStore()
	body, err, point := nodeDefault(body, bacnetAI, category, store.GetApplication())
	var pointUUID string
	if point != nil {
		pointUUID = point.UUID
	}
	return &AI{
		body,
		false,
		0,
		points.AnalogInput,
		pointUUID,
	}, err
}

func (inst *AI) subscribePresentValue() {
	//topicPv := TopicPresentValue(typeAI, inst.objectID)
	//getMqtt().Subscribe(topicPv)
}

func (inst *AI) setObjectId() {
	id, ok := conversions.GetInt(inst.ReadPin(node.ObjectId))
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}

func (inst *AI) getObjectId() (int, bool) {
	return conversions.GetInt(inst.ReadPin(node.ObjectId))
}

func (inst *AI) Process() {
	if !inst.onStart {
		inst.setObjectId()
	}
	updateInputs(inst, points.AnalogInput, inst.objectID)
	inst.onStart = true
}

func (inst *AI) Cleanup() {}
