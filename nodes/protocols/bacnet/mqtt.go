package bacnet

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/topics"
	points2 "github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	eventbus2 "github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func getTopic(msg interface{}) string {
	m := decode(msg)
	if m != nil {
		return m.Msg.Topic()
	}
	return ""
}

func decode(msg interface{}) *eventbus2.Message {
	m, ok := msg.(*eventbus2.Message)
	if ok {
		return m
	}
	return nil
}

var priorityBus runnerStatus

func (inst *Server) priorityBus() {
	if !priorityBus {
		handlerMQTT := bus.Handler{
			Handle: func(ctx context.Context, e bus.Event) {
				go func() {
					decoded := decode(e.Data)
					if decoded != nil {

					}

				}()
			},
			Matcher: eventbus2.BacnetPri,
		}
		key := fmt.Sprintf("key_%s", helpers.UUID())
		eventbus2.GetBus().RegisterHandler(key, handlerMQTT)
	}
	priorityBus = true
}

func (inst *Server) presetValueBus() {
	if !priorityBus {
		handlerMQTT := bus.Handler{
			Handle: func(ctx context.Context, e bus.Event) {
				go func() {
					decoded := decode(e.Data)
					if decoded != nil {

					}
				}()
			},
			Matcher: eventbus2.BacnetPV,
		}
		key := fmt.Sprintf("key_%s", helpers.UUID())
		eventbus2.GetBus().RegisterHandler(key, handlerMQTT)
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
			Matcher: eventbus2.RubixIOInputs,
		}
		key := fmt.Sprintf("key_%s", helpers.UUID())
		eventbus2.GetBus().RegisterHandler(key, handlerMQTT)
	}
	priorityBus = true
}

//mqttRunner processMessage are the messages from the bacnet-server via the mqtt-broker
func (inst *Server) mqttSubRunner() {
	log.Info("start mqtt-sub-runner")
}

func (inst *Server) handleBacnet(msg interface{}) {
	payload := points2.NewPayload()
	err := payload.NewMessage(msg)
	if err != nil {
		log.Errorf("bacnet-sub-runner malformed mqtt message err:%s", err.Error())
		return
	}
	topic := payload.GetTopic()
	t, id := payload.GetObjectID()
	point := getStore().GetPointByObject(t, id)
	if topics.IsPri(topic) {
		value := payload.GetHighestPriority()
		log.Infof("mqtt-runner-subscribe point type:%s-%d value:%f", point.ObjectType, point.ObjectID, value.Value)
		getStore().WritePointValue(point.UUID, value.Value)
	}
}

//mqttPubRunner send messages to the broker, as in read a modbus point and send it to the bacnet server
func (inst *Server) mqttPubRunner() {
	log.Info("start mqtt-pub-runner")
	for {
		for _, point := range getStore().GetPoints() {
			if point.ToBacnetSyncPending {
				inst.mqttPublish(point)
				getStore().UpdateBacnetSync(point.UUID, false)
			} else {
				//log.Infof("mqtt-runner-publish point skip as non cov type:%s-%d", point.ObjectType, point.ObjectID)
			}
		}
		time.Sleep(runnerDelay * time.Millisecond)
	}
}

func (inst *Server) mqttPublish(pnt *points2.Point) {
	if pnt == nil {
		log.Errorf("bacnet-server-publish point can not be empty")
		return
	}
	objectType := pnt.ObjectType
	objectId := pnt.ObjectID
	value := pnt.ToBacnet
	obj, err := points2.ObjectSwitcher(objectType)
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
