package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/mqttclient"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

type MqttPub struct {
	*node.BaseNode
	client     *mqttclient.Client
	connected  bool
	subscribed bool
	newMessage string
	mqttTopic  string
}

func NewMqttPub(body *node.BaseNode) (node.Node, error) {
	body = node.Defaults(body, mqttPublish, category)
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
		mqttTopic = ""
	}
	inputs := node.BuildInputs(node.BuildInput(node.In1, node.TypeString, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return &MqttPub{body, nil, false, false, "", mqttTopic}, nil
}

func (inst *MqttPub) getTopic() string {
	return inst.mqttTopic
}

func (inst *MqttPub) publish(value string) {
	c, _ := mqttclient.GetMQTT()
	if inst.getTopic() != "" {
		err := c.Publish(inst.getTopic(), mqttclient.AtMostOnce, true, value)
		if err != nil {
			log.Errorf(fmt.Sprintf("mqtt-publish topic:%s err:%s", inst.getTopic(), err.Error()))
		}
	} else {
		log.Errorf(fmt.Sprintf("mqtt-publish topic can not be empty"))
	}
}

func (inst *MqttPub) connect() {
	mqttBroker := "tcp://0.0.0.0:1883"
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		log.Errorf(fmt.Sprintf("mqtt-publish-connect err:%s", err.Error()))
	}
	client, connected := mqttclient.GetMQTT()
	inst.connected = connected
	inst.client = client
}

func (inst *MqttPub) Process() {
	in1Pointer, in1Val := inst.ReadPin(node.In1)
	if in1Pointer == nil {
		return
	}
	if !inst.connected {
		go inst.connect()
	}
	if inst.connected {
		go inst.publish(in1Val)
	}
	val := float.StrToFloat(in1Val)
	inst.WritePin(node.Out1, float.ToStrPtr(val))

}

func (inst *MqttPub) Cleanup() {}
