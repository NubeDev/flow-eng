package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
	log "github.com/sirupsen/logrus"
)

type BV struct {
	*node.Spec
	connected  bool
	subscribed bool
	objectID   bstore.ObjectID
	objectType bstore.ObjectType
	pointUUID  string
}

const (
	object = "object"
)

var objects = []string{"analog_value", "binary_value"}

func NewBacnetBV(body *node.Spec) (node.Node, error) {
	var err error
	store := getStore()
	body, err, point := nodeDefault(body, bacnetBV, category, store.GetApplication())
	if err != nil {
		log.Error(err)
		//return nil, err
	}
	var pointUUID string
	if point != nil {
		pointUUID = point.UUID
	}

	return &BV{
		body,
		false,
		false,
		0,
		bstore.BinaryVariable,
		pointUUID,
	}, nil
}

func (inst *BV) subscribePresentValue() {
	topicPv := TopicPresentValue(typeBV, inst.objectID)
	inst.client().Subscribe(topicPv)
}

func (inst *BV) subscribePriority() {
	topicPriority := TopicPriority(typeBV, inst.objectID)
	inst.client().Subscribe(topicPriority)
}

func (inst *BV) client() *mqttbase.Mqtt {
	return client
}

func (inst *BV) bus() cbus.Bus {
	return inst.client().BACnetBus()
}

func (inst *BV) setObjectId() {
	id, ok := inst.ReadPin(node.ObjectId).(int)
	if ok {
		inst.objectID = bstore.ObjectID(id)
	}
}

func (inst *BV) setConnected() {
	inst.connected = true
}

func (inst *BV) setDisconnected() {
	inst.connected = false
}

var loopCount uint64

func (inst *BV) Process() {
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

func (inst *BV) Cleanup() {}
