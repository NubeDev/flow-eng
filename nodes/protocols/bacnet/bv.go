package bacnet

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
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
	id, ok := inst.ReadPin(node.ObjectId).(int)
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}

var loopCount uint64

func (inst *BV) Process() {
	loopCount++
	//if !getMqtt().Connected() || !inst.connected {
	//	inst.setObjectId()
	//	inst.subscribePriority()
	//	inst.connected = true
	//}
	//if !getMqtt().Connected() {
	//	inst.connected = false
	//}

}

func (inst *BV) Cleanup() {}
