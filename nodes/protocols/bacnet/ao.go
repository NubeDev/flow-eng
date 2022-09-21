package bacnet

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

type AO struct {
	*node.Spec
	connected  bool
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
	id, ok := inst.ReadPin(node.ObjectId).(int)
	if ok {
		inst.objectID = points.ObjectID(id)
	}
}

func (inst *AO) Process() {
	process(inst)

	//if !getMqtt().Connected() || !inst.connected {
	//	inst.setObjectId()
	//	inst.connected = true
	//}
	//if !getMqtt().Connected() {
	//	inst.connected = false
	//}

}

func (inst *AO) Cleanup() {}
