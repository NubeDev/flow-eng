package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/mqttclient"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schema"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MqttSub struct {
	*node.BaseNode
	client     *mqttclient.Client
	connected  bool
	subscribed bool
	newMessage string
}

const (
	topic = "topic"
)

var bus cbus.Bus

func NewMqttSub(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body, mqttSub)
	body.Info.Name = mqttSub
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeString, nil, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeString, body.Outputs))
	decode := schema.NewString(nil)
	err := body.DecodeProperties(topic, decode)
	if err != nil {
		return nil, err
	}
	t := schema.NewString(&schema.SettingBase{
		Title:        topic,
		Min:          1,
		DefaultValue: decode.DefaultValue,
	})
	settings, err := node.BuildSettings(node.BuildSetting(schema.PropString, topic, t))
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

func (inst *MqttSub) getTopic() string {
	val, err := inst.GetPropValueStr(topic)
	if err != nil {
		return val
	}
	return val

}

func (inst *MqttSub) subscribe() {
	c, _ := mqttclient.GetMQTT()
	if inst.getTopic() != "" {
		err := c.Subscribe(inst.getTopic(), mqttclient.AtMostOnce, handle)
		if err != nil {
			log.Errorf(fmt.Sprintf("mqtt-subscribe topic:%s err:%s", inst.getTopic(), err.Error()))
		}
		inst.subscribed = true
	} else {
		log.Errorf(fmt.Sprintf("mqtt-subscribe topic can not be empty"))
	}
}

func (inst *MqttSub) connect() {
	mqttBroker := "tcp://0.0.0.0:1883"
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		log.Errorf(fmt.Sprintf("mqtt-subscribe-connect err:%s", err.Error()))
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
			log.Info("MQTT:newMessage", inst.newMessage)
		}
	}()

	val := float.StrToFloat(inst.newMessage)

	inst.WritePin(node.Out1, float.ToStrPtr(val))

}

func (inst *MqttSub) Cleanup() {}
