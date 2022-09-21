package eventbus

import (
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/topics"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	UUID string
	Msg  mqtt.Message
}

var PointsHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	mes := &Message{helpers.ShortUUID("bus"), msg}
	topic := msg.Topic()
	if topics.CheckBACnet(topic) {
		if topics.IsPV(topic) {
			GetService().RegisterTopic(BacnetPV)
			err := GetService().Emit(CTX(), BacnetPV, mes)
			if err != nil {
				return
			}
		}
		if topics.IsPri(topic) {
			GetService().RegisterTopic(BacnetPri)
			err := GetService().Emit(CTX(), BacnetPri, mes)
			if err != nil {
				return
			}
		}
	}
	if topics.CheckRubixIO(topic) {
		GetService().RegisterTopic(RubixIOInputs)
		err := GetService().Emit(CTX(), RubixIOInputs, mes)
		if err != nil {
			return
		}
	}
}
