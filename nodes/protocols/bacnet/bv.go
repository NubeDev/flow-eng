package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/points"
)

type BV struct {
	*node.Spec
	connected  bool
	objectID   points.ObjectID
	objectType points.ObjectType
	pointUUID  string
}

const (
	object = "object"
)

func NewBV(body *node.Spec) (node.Node, error) {
	var err error
	store := getStore()
	body, err, point := nodeDefault(body, bacnetBV, category, store.GetApplication())
	var pointUUID string
	if point != nil {
		pointUUID = point.UUID
	}
	return &BV{
		body,
		false,
		0,
		points.BinaryVariable,
		pointUUID,
	}, err
}

func (inst *BV) subscribePriority() {
	topicPriority := TopicPriority(typeBV, inst.objectID)
	getClient().Subscribe(topicPriority)
}

func (inst *BV) bus() cbus.Bus {
	return getClient().BACnetBus()
}

func (inst *BV) setObjectId() {
	id, ok := inst.ReadPin(node.ObjectId).(int)
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}

var loopCount uint64

func (inst *BV) Process() {
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

func (inst *BV) Cleanup() {}
