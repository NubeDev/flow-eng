package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/mqttclient"
	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
)

type MqttSub struct {
	*node.BaseNode
	client     *mqttclient.Client
	connected  bool
	subscribed bool
	newMessage string
}

var bus cbus.Bus

type topic struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"topic"`
	Min      int    `json:"minLength" default:"1"`
	ReadOnly bool   `json:"readOnly" default:"false"`
	Value    string `json:"value"`
}

func NewMqttSub(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body)
	body.Info.Name = mqttSub
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeString, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeString, body.Outputs))

	topicSetting := &topic{}
	aaa := body.GetProperties("topic")
	err := mapstructure.Decode(aaa, topicSetting)
	if err != nil {
		return nil, err
	}

	settings, err := node.BuildSettings(node.BuildSetting("string", "topic", &topic{
		Type:     "string",
		Title:    "topic",
		Min:      1,
		ReadOnly: true,
		Value:    topicSetting.Value,
	}), node.BuildSetting("string", "topic2", &topic{
		Type:     "string",
		Title:    "topic",
		Min:      1,
		ReadOnly: true,
		Value:    "",
	}))
	if err != nil {
		return nil, err
	}

	body.Settings = settings

	bus = cbus.New(1)
	return &MqttSub{body, nil, false, false, ""}, nil
}

var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	bus.Send(msg)
}

func (inst *MqttSub) subscribe() {
	topicSetting := &topic{}
	topicSetting, ok := inst.GetProperties("topic").(*topic)
	if ok {
		c, _ := mqttclient.GetMQTT()
		c.Subscribe(topicSetting.Value, mqttclient.AtMostOnce, handle)
		inst.subscribed = true
	}

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

	val := float.StrToFloat(inst.newMessage)

	inst.WritePin(node.Out1, float.ToStrPtr(val))
	fmt.Println("****************newMessage", inst.newMessage)
}

func (inst *MqttSub) Cleanup() {}
