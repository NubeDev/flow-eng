package bacnet

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/mustafaturan/bus/v3"
)

var priorityBus runnerStatus

func (inst *Server) priorityBus() {
	if !priorityBus {
		handlerMQTT := bus.Handler{
			Handle: func(ctx context.Context, e bus.Event) {
				go func() {
					decoded := decode(e.Data)
					if decoded != nil {
						inst.fromBacnet(decoded) // this messages will come from 3rd party bacnet devices
					}
				}()
			},
			Matcher: eventbus.BacnetPri,
		}
		key := fmt.Sprintf("key_%s", helpers.UUID())
		eventbus.GetBus().RegisterHandler(key, handlerMQTT)
	}
	priorityBus = true
}
