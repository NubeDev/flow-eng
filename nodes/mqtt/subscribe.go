package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/mqttclient"
	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MqttSub struct {
	*node.Spec
	client     *mqttclient.Client
	connected  bool
	subscribed bool
	newMessage string
	mqttTopic  string
}

const (
	topic = "topic"
)

var bus cbus.Bus

func NewMqttSub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, mqttSub, category)
	_, setting, value, err := node.NewSetting(body, &node.SettingOptions{Type: node.String, Title: topic, Min: 1, Max: 200})
	if err != nil {
		return nil, err
	}
	settings, err := node.BuildSettings(setting)
	if err != nil {
		return nil, err
	}
	mqttTopic, ok := value.(string)
	if !ok {
		mqttTopic = "sub"
	}
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	bus = cbus.New(1)
	return &MqttSub{body, nil, false, false, "", mqttTopic}, nil
}

var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	bus.Send(msg)
}

func (inst *MqttSub) getTopic() string {
	return inst.mqttTopic
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
	if bus == nil {
		panic("mqtt-bus can not be empty")
	}
	if !inst.connected {
		go inst.connect()
	}
	if inst.connected {
		if !inst.subscribed {
			inst.subscribe()
		}
	}
	if inst.connected && inst.subscribed {
		go func() {
			msg, ok := bus.Recv()
			if ok {
				msg_ := msg.(mqtt.Message)
				inst.newMessage = string(msg_.Payload())
				log.Info("MQTT:newMessage", inst.newMessage)
			}
		}()

		val := float.StrToFloat(inst.newMessage)
		inst.WritePin(node.Out, float.ToStrPtr(val))
	}
}

func (inst *MqttSub) Cleanup() {}
