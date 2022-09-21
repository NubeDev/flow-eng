package bacnet

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
	"time"
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

// mqttPubRunner send messages to the broker, as in read a modbus point and send it to the bacnet server
func (inst *Server) writeRunner() {
	log.Info("start mqtt-pub-runner")
	for {
		for _, point := range getStore().GetPoints() {
			if point.WriteSyncPending {
				inst.mqttPublish(point)
				getStore().UpdateBacnetSync(point.UUID, false)
			} else {
				//log.Infof("mqtt-runner-publish point skip as non cov type:%s-%d", point.ObjectType, point.ObjectID)
			}
		}
		time.Sleep(runnerDelay * time.Millisecond)
	}
}

// mqttPublish example for future MQTT write to the bacnet-server
func (inst *Server) mqttPublish(pnt *points.Point) {
	if pnt == nil {
		log.Errorf("bacnet-server-publish point can not be empty")
		return
	}
	objectType := pnt.ObjectType
	objectId := pnt.ObjectID
	value := pnt.WriteValue
	obj, err := points.ObjectSwitcher(objectType)
	if err != nil {
		log.Error(err)
		return
	}
	topic := fmt.Sprintf("bacnet/%s/%d", obj, objectId)
	err = inst.client.Publish(topic, mqttclient.AtMostOnce, true, value)
	if err != nil {
		return
	}
}

func getTopic(msg interface{}) string {
	m := decode(msg)
	if m != nil {
		return m.Msg.Topic()
	}
	return ""
}

func decode(msg interface{}) *eventbus.Message {
	m, ok := msg.(*eventbus.Message)
	if ok {
		return m
	}
	return nil
}
