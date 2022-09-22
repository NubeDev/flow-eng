package bacnet

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/NubeDev/flow-eng/services/rubixio"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (inst *Server) rubixIOBus() {
	if !priorityBus {
		handlerMQTT := bus.Handler{
			Handle: func(ctx context.Context, e bus.Event) {
				go func() {
					decoded := decode(e.Data)
					if decoded != nil {
						inst.rubixInputsRunner(decoded)
					}
				}()
			},
			Matcher: eventbus.RubixIOInputs,
		}
		key := fmt.Sprintf("key_%s", helpers.UUID())
		eventbus.GetBus().RegisterHandler(key, handlerMQTT)
	}
	priorityBus = true
}

func (inst *Server) rubixInputsRunner(msg *eventbus.Message) {
	rubix := &rubixIO.RubixIO{}
	inputs, err := rubix.DecodeInputs(msg.Msg.Payload())
	if err != nil {
		log.Error(err)
		//return
	}

	for _, point := range getStore().GetPointsByApplication(applications.RubixIO) {
		if point.ObjectType == points.AnalogInput {
			value, err := rubix.GetInputValue(point, inputs)
			if err != nil {
				return
			}
			getStore().WriteValueFromRead(point.UUID, value)
			fmt.Println(value, "rubix-io-input-value")
		}
	}

}
