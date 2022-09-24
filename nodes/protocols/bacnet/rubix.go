package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/rubixio"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var rubixIOBus mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mes := &topics.Message{UUID: helpers.ShortUUID("bus"), Msg: msg}
	if topics.CheckRubixIO(msg.Topic()) {
		inst.rubixInputsRunner(mes)
	}

}

func (inst *Server) rubixInputsRunner(msg *topics.Message) {
	rubix := &rubixIO.RubixIO{}
	inputsPayload, err := rubix.DecodeInputs(msg.Msg.Payload()) // data comes from mqtt
	if err != nil {
		log.Error(err)
		//return
	}

	for _, point := range getStore().GetPointsByApplication(applications.RubixIO) {
		if point.ObjectType == points.AnalogInput {
			value, err := rubix.DecodeInputValue(point, inputsPayload)
			if err != nil {
				return
			}
			getStore().WriteValueFromRead(point.UUID, value)
			fmt.Println(value, "rubix-io-input-value")
		}
	}

}
