package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/topics"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var bacnetBus mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mes := &topics.Message{UUID: helpers.ShortUUID("bus"), Msg: msg}
	if topics.IsPri(msg.Topic()) {
		inst.fromBacnet(mes)
	}

}
