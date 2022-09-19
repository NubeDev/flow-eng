package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

type AI struct {
	*node.Spec
	connected  bool
	subscribed bool
	objectID   bstore.ObjectID
	objectType bstore.ObjectType
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
		false,
		0,
		bstore.AnalogInput,
		pointUUID,
	}, err
}

func (inst *AI) subscribePresentValue() {
	topicPv := TopicPresentValue(typeAI, inst.objectID)
	inst.client().Subscribe(topicPv)
}

func (inst *AI) subscribePriority() {
	topicPriority := TopicPriority(typeAI, inst.objectID)
	inst.client().Subscribe(topicPriority)
}

func (inst *AI) client() *mqttbase.Mqtt {
	return client
}

func (inst *AI) bus() cbus.Bus {
	return inst.client().BACnetBus()
}

func (inst *AI) setObjectId() {
	id, ok := inst.ReadPin(node.ObjectId).(int)
	if ok {
		inst.objectID = bstore.ObjectID(id)
	}
}

func (inst *AI) setConnected() {
	inst.connected = true
}

func (inst *AI) setDisconnected() {
	inst.connected = false
}

func (inst *AI) Process() {
	loopCount++
	if !inst.connected {
		inst.setObjectId()
		inst.client().Connected()
		inst.setConnected()
		inst.subscribePresentValue()

	}

	if inst.connected {
		//inst.processMessage()
	}

}

func (inst *AI) Cleanup() {}
