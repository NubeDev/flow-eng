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
}

func NewMqttPub(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body, mqttPublish)
	body.Info.Name = node.SetName(mqttPublish)
	body.Info.Category = node.SetName(category)
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)

	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeString, nil, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeString, nil, body.Outputs))

	//setting, err := node.BuildSetting(schema.PropString, topic, body)
	//if err != nil {
	//	return nil, err
	//}
	//settings, err := node.BuildSettings(setting)
	//if err != nil {
	//	return nil, err
	//}
	//body.Settings = settings

	return &MqttPub{body, nil, false, false, ""}, nil
}

func (inst *MqttPub) getTopic() string {
	val, err := inst.GetPropValueStr(topic)
	if err != nil {
		return val
	}
	return val
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
