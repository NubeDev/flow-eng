package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

type Read struct {
	*node.Spec
	client         *mqttbase.Mqtt
	connected      bool
	subscribed     bool
	latestPv       *float64
	latestPriority *PriArray
	topicPv        string
	topicPriority  string
}

const (
	object = "object"
)

var objects = []string{"analog_value", "binary_value"}

// outputs are both pub/sub to the server
// inputs don't need to subscribe only publish to the server

// bacnet/program/2508/state -> "started" on start of the bacnet-server
// bacnet/ao/1/pv -> "11.000000"
// bacnet/ao/1/pri -> {Null,33.3,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,11.000000}

func NewBacnetBVRead(body *node.Spec, opts interface{}) (node.Node, error) {
	var err error
	var client *mqttbase.Mqtt
	var point *bstore.Point
	//store := GetStore()
	body, client, err, point = nodeDefault(body, bacnetReadBV, category, applications.BACnet, opts)
	if err != nil {
		return nil, err
	}
	if client == nil {
	}
	fmt.Println(point)

	//err = GetStore().AddPoint(point)

	bacnet := &Read{body, client, false, false, nil, nil, "", ""}
	return bacnet, nil
}

func (inst *Read) subscribePresentValue(address float64) {
	//inst.topicPv = TopicPresentValue(bv, GetObjectId(address))
	// bacnet/ao/1/pv
	inst.client.Subscribe("bacnet/ao/1/pv")
}

func (inst *Read) subscribePriority(address float64) {

	inst.topicPriority = TopicPriority(typeBV, GetObjectId(address))
	inst.client.Subscribe(inst.topicPriority)
}

func (inst *Read) bus() cbus.Bus {
	return inst.client.BACnetBus()
}

func (inst *Read) processMessage() {
	go func() {
		msg, ok := inst.bus().Recv()
		if ok {

			m, ok := msg.(*mqttbase.Message)
			if ok {
				fmt.Println(m.Msg.Topic(), string(m.Msg.Payload()))
			}

			//log.Info("MQTT:newMessage", inst.newMessage)
		}
	}()

}

func (inst *Read) setConnected() {
	inst.connected = true
}

func (inst *Read) setDisconnected() {
	inst.connected = false
}

var loopCount uint64

func (inst *Read) Process() {
	loopCount++
	if !inst.connected {
		inst.client.Connected()
		inst.setConnected()
		inst.subscribePresentValue(1)
	}

	if inst.connected {

		inst.processMessage()
	}

}

func (inst *Read) Cleanup() {}
