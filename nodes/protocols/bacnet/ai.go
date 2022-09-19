package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/points"
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
	topicPv := TopicPresentValue(typeAI, inst.objectID)
	getClient().Subscribe(topicPv)
}

func (inst *AI) subscribePriority() {
	topicPriority := TopicPriority(typeAI, inst.objectID)
	getClient().Subscribe(topicPriority)
}

func (inst *AI) bus() cbus.Bus {
	return getClient().BACnetBus()
}

func (inst *AI) setObjectId() {
	id, ok := inst.ReadPin(node.ObjectId).(int)
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}

func (inst *AI) Process() {
	loopCount++
	if !getClient().Connected() || !inst.connected {
		inst.setObjectId()
		inst.subscribePriority()
		inst.connected = true
	}
	if !getClient().Connected() {
		inst.connected = false
	}

}

func (inst *AI) Cleanup() {}
