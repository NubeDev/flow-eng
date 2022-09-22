package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

type AI struct {
	*node.Spec
	connected  bool
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

	id, ok := getInt(inst.ReadPin(node.ObjectId))
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}

func (inst *AI) getObjectId() (int, bool) {
	return getInt(inst.ReadPin(node.ObjectId))

}

func (inst *AI) Process() {
	id, _ := inst.getObjectId()
	fmt.Println("ID", id)

	v, _ := getStore().GetValueFromReadByObject(points.AnalogInput, 1)
	fmt.Println("VALUE", v)
	inst.WritePin(node.Out, v)

}

func (inst *AI) Cleanup() {}
