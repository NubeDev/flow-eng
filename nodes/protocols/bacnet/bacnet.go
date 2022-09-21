package bacnet

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

var priorityBus runnerStatus

func (inst *Server) priorityBus() {
	if !priorityBus {
		handlerMQTT := bus.Handler{
			Handle: func(ctx context.Context, e bus.Event) {
				go func() {
					decoded := decode(e.Data)
					if decoded != nil {
						inst.handleBacnet(decoded) // this messages will come from 3rd party bacnet devices
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

func (inst *Server) setToSync() points.SyncTo {
	app := getRunnerType()
	switch app {
	case applications.RubixIO:
		return points.ToRubixIO
	}
	return ""

}

func (inst *Server) handleBacnet(msg interface{}) {
	payload := points.NewPayload()
	err := payload.NewMessage(msg)
	if err != nil {
		log.Errorf("bacnet-sub-runner malformed mqtt message err:%s", err.Error())
		return
	}
	topic := payload.GetTopic()
	t, id := payload.GetObjectID()
	point := getStore().GetPointByObject(t, id)
	if point == nil {
		log.Errorf("mqtt-payload-priorty-array no point-found in store for type:%s-%d", t, id)
		return
	}
	if topics.IsPri(topic) {
		value := payload.GetHighestPriority()
		log.Infof("mqtt-runner-subscribe point type:%s-%d value:%f", point.ObjectType, point.ObjectID, value.Value)
		getStore().AddSync(point.UUID, value.Value, points.FromMqttPriory, inst.setToSync())
		getStore().WritePointValue(point.UUID, value.Value)
	}

}
