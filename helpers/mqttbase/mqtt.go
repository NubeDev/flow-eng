package mqttbase

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Mqtt struct {
	client     *mqttclient.Client
	connected  bool
	subscribed bool
}

var mqttBus cbus.Bus
var bacnetBus cbus.Bus

func NewMqtt() (*Mqtt, error) {
	mqttBus = cbus.New(1)
	bacnetBus = cbus.New(1)
	return &Mqtt{}, nil
}

type Message struct {
	UUID string
	Msg  mqtt.Message
}

func checkBACnet(topic string) (isBacnet bool) { // to try and save spamming random message
	parts := strings.Split(topic, "/")
	if len(parts) > 0 {
		if parts[0] == "bacnet" {
			return true
		}
	}
	return isBacnet
}

var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	if checkBACnet(msg.Topic()) {
		mes := &Message{helpers.ShortUUID("bac"), msg}
		bacnetBus.Send(mes)
	}
}

func (inst *Mqtt) BACnetBus() cbus.Bus {
	return bacnetBus
}

func (inst *Mqtt) getClient() *mqttclient.Client {
	return inst.client
}

func (inst *Mqtt) Publish(value interface{}, topic string) {
	c := inst.getClient()
	if topic != "" {
		v := fmt.Sprintf("%v", value)
		err := c.Publish(topic, mqttclient.AtMostOnce, true, v)
		if err != nil {
			log.Errorf(fmt.Sprintf("mqttbase-publish topic:%s err:%s", topic, err.Error()))
		}
	} else {
		log.Errorf(fmt.Sprintf("mqttbase-publish topic can not be empty"))
	}
}

func (inst *Mqtt) Subscribe(topic string) {
	c := inst.getClient()
	if topic != "" {
		err := c.Subscribe(topic, mqttclient.AtMostOnce, handle)
		if err != nil {
			log.Errorf(fmt.Sprintf("mqttbase-subscribe topic:%s err:%s", topic, err.Error()))
		}
		inst.subscribed = true
	} else {
		log.Errorf(fmt.Sprintf("mqttbase-subscribe topic can not be empty"))
	}
}

func (inst *Mqtt) Connected() bool {
	return inst.connected
}

func (inst *Mqtt) Connect() {
	mqttBroker := "tcp://0.0.0.0:1883"
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		log.Errorf(fmt.Sprintf("mqttbase-subscribe-connect err:%s", err.Error()))
	}
	client, connected := mqttclient.GetMQTT()
	inst.connected = connected
	inst.client = client
}
