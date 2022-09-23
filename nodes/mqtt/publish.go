package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
)

type MqttPub struct {
	*node.Spec
	client     *mqttclient.Client
	connected  bool
	subscribed bool
	newMessage string
	mqttTopic  string
}

func NewMqttPub(body *node.Spec) (node.Node, error) {
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
	inputs := node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
	return &MqttPub{body, nil, false, false, "", mqttTopic}, nil
}

func (inst *MqttPub) getTopic() string {
	str, err := inst.GetPropValueStr(topic)
	if err != nil {
		return ""
	}
	inst.mqttTopic = str
	return inst.mqttTopic
}

func (inst *MqttPub) publish(value interface{}) {
	c, _ := mqttclient.GetMQTT()
	if inst.getTopic() != "" {
		v := fmt.Sprintf("%v", value)
		err := c.Publish(inst.getTopic(), mqttclient.AtMostOnce, true, v)
		if err != nil {
			log.Errorf(fmt.Sprintf("pointbus-publish topic:%s err:%s", inst.getTopic(), err.Error()))
		}
	} else {
		log.Errorf(fmt.Sprintf("pointbus-publish topic can not be empty"))
	}
}

func (inst *MqttPub) connect() {
	mqttBroker := "tcp://0.0.0.0:1883"
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		log.Errorf(fmt.Sprintf("pointbus-publish-connect err:%s", err.Error()))
	}
	client, connected := mqttclient.GetMQTT()
	inst.connected = connected
	inst.client = client
}

func (inst *MqttPub) Process() {
	val := inst.ReadPin(node.In1)
	if val == nil {
		return
	}
	if !inst.connected {
		go inst.connect()
	}
	if inst.connected {
		go inst.publish(val)
	}
	inst.WritePin(node.Out1, val)
}

func (inst *MqttPub) Cleanup() {}
