package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/mqttclient"
	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttSub struct {
	*node.BaseNode
	client     *mqttclient.Client
	connected  bool
	subscribed bool
	newMessage string
}

var bus cbus.Bus

func NewMqttSub(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body)
	body.Info.Name = mqttSub
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.Topic, node.TypeString, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeString, body.Outputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out2, node.TypeFloat, body.Outputs))
	bus = cbus.New(1)
	return &MqttSub{body, nil, false, false, ""}, nil
}

var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	bus.Send(msg)
}

func (inst *MqttSub) subscribe() {
	c, _ := mqttclient.GetMQTT()
	c.Subscribe("hello", mqttclient.AtMostOnce, handle)
	inst.subscribed = true
}

func (inst *MqttSub) connect() {
	mqttBroker := "tcp://0.0.0.0:1883"
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		return
	}
	client, connected := mqttclient.GetMQTT()
	inst.connected = connected
	inst.client = client
}

func (inst *MqttSub) Process() {
	if !inst.connected {
		go inst.connect()
	}
	if inst.connected {
		if !inst.subscribed {
			inst.subscribe()
		}
	}

	go func() {
		msg, ok := bus.Recv()
		if ok {
			msg_ := msg.(mqtt.Message)
			inst.newMessage = string(msg_.Payload())
		}
	}()
	fmt.Println("newMessage", inst.newMessage)
}

func (inst *MqttSub) Cleanup() {}
