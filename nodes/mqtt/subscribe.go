package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MqttSub struct {
	*node.Spec
	client     *mqttclient.Client
	connected  bool
	subscribed bool
}

const (
	topic = "topic"
)

var instMqttSub *MqttSub
var newMessage string

func NewMqttSub(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, mqttSub, category)
	//_, setting, _, err := node.NewSetting(body, &node.SettingOptions{Type: node.String, Title: topic, Min: 1, Max: 200})
	//if err != nil {
	//	return nil, err
	//}
	//settings, err := node.BuildSettings(setting)
	//if err != nil {
	//	return nil, err
	//}

	top := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	ip := node.BuildInput(node.Ip, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(top, ip)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	instMqttSub = &MqttSub{body, nil, false, false}
	return instMqttSub, nil
}

var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	//fmt.Println(string(msg.Payload()))
	newMessage = string(msg.Payload())

}

func getMqttSub() *MqttSub {
	return instMqttSub
}

func (inst *MqttSub) subscribe(topic string) {
	if topic != "" {
		err := inst.client.Subscribe(topic, mqttclient.AtMostOnce, handle)
		if err != nil {
			log.Errorf(fmt.Sprintf("mqtt-subscribe topic:%s err:%s", topic, err.Error()))
		}
		inst.subscribed = true
	} else {
		log.Errorf(fmt.Sprintf("mqtt-subscribe topic can not be empty"))
	}
}

func (inst *MqttSub) connect(ip string) {
	var err error
	fmt.Println("MQTT CONNECT")
	if ip == "" {
		ip = "0.0.0.0"
	}
	inst.client, err = mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{fmt.Sprintf("tcp://%s:1883", ip)},
	})
	err = inst.client.Connect()
	if err != nil {
		log.Error(err)
		//return nil, err
	}
	inst.connected = true
}

func (inst *MqttSub) Process() {
	top := inst.ReadPinAsString(node.Topic)
	ip := inst.ReadPinAsString(node.Ip)

	if !inst.connected {
		go inst.connect(ip)
	}
	if inst.connected {
		if !inst.subscribed {
			go inst.subscribe(top)
		}
	}
	if newMessage != "" {
		inst.WritePin(node.Out, newMessage)
	} else {
		inst.WritePin(node.Out, nil)
	}
}

func (inst *MqttSub) Cleanup() {}
