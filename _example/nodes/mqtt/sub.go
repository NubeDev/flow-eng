package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/mqttclient"
	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type Mqtt struct {
	*node.BaseNode
	client     *mqttclient.Client
	connected  bool
	subscribed bool
}

func NewMqtt(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body)
	body.Info.Name = mqttNode
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, body.Inputs), node.BuildInput(node.In2, node.TypeFloat, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, body.Outputs))
	return &Mqtt{body, nil, false, false}, nil
}

// used for getting data into the plugins
var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
}

func (inst *Mqtt) subscribe() {
	c, _ := mqttclient.GetMQTT()
	c.Subscribe("hello", mqttclient.AtMostOnce, handle)
	inst.subscribed = true
}

func (inst *Mqtt) connect() {
	mqttBroker := "tcp://0.0.0.0:1883"
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		return
	}
	client, connected := mqttclient.GetMQTT()
	fmt.Println("MQTT client connect****************", connected)
	inst.connected = connected
	inst.client = client
}

func (inst *Mqtt) Process() {
	if !inst.connected {
		go inst.connect()
	}
	if inst.connected {
		if !inst.subscribed {
			inst.subscribe()
		}
	}

}

func (inst *Mqtt) Cleanup() {}
