package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type Broker struct {
	*node.Spec
	firstLoop     bool
	loopCount     uint64
	connection    *db.Connection
	mqttClient    *mqttclient.Client
	mqttConnected bool
}

var mqttQOS = mqttclient.AtMostOnce
var mqttRetain = false

func NewBroker(body *node.Spec) (node.Node, error) {
	//var err error
	body = node.Defaults(body, mqttBroker, category)
	inputs := node.BuildInputs()
	outputs := node.BuildOutputs(node.BuildOutput(node.Connected, node.TypeBool, nil, body.Outputs))
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	network := &Broker{body, false, 0, nil, nil, false}
	//
	network.Spec = body
	return network, nil
}

func getPayloads(children interface{}, ok bool) []*mqttPayload {
	if ok {
		mqttData := children.(*mqttStore)
		if mqttData != nil {
			return mqttData.payloads
		}
	}
	return nil
}

func (inst *Broker) subscribe() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		s := inst.GetStore()
		children, ok := s.Get(inst.GetID())
		payloads := getPayloads(children, ok)
		for _, payload := range payloads {
			if payload.topic == message.Topic() {
				n := inst.GetNode(payload.nodeUUID)
				n.SetPayload(&node.Payload{
					Any: string(message.Payload()),
				})
			}
		}
	}
	s := inst.GetStore()
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	for _, payload := range payloads {
		if payload.topic != "" {
			if inst.mqttClient != nil {
				err := inst.mqttClient.Subscribe(payload.topic, mqttQOS, callback)
				if err != nil {
					log.Errorf("mqtt-broker subscribe:%s err:%s", payload.topic, err.Error())
				} else {
					log.Infof("mqtt-broker subscribe:%s", payload.topic)
				}
			}
		}
	}
}

func (inst *Broker) publish() {
	s := inst.GetStore()
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	for _, payload := range payloads {
		if !payload.isPublisher {
			continue
		}
		n := inst.GetNode(payload.nodeUUID)
		t, _ := n.ReadPinAsString(node.Topic)
		p, _ := n.ReadPinAsString(node.Message)
		if p != "" {
			if t == payload.topic {
				if inst.mqttClient != nil {
					updated, _ := n.InputUpdated(node.Message)
					if updated {
						err := inst.mqttClient.Publish(t, mqttQOS, mqttRetain, p)
						if err != nil {
							log.Errorf("mqtt-broker publish err:%s", err.Error())
						}
					}

				}
			}
		}
	}
}

func (inst *Broker) setConnection() {
	settings, err := getSettings(inst.GetSettings())
	if err != nil {
		log.Errorf("add mqtt broker failed to get settings err:%s", err.Error())
		return
	}
	connection, err := inst.GetDB().GetConnection(settings.Conn)
	if err != nil {
		log.Error("add mqtt broker failed to find connection")
		return
	}
	inst.connection = connection
	mqttClient, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{fmt.Sprintf("tcp://%s:%d", connection.Host, connection.Port)},
	})
	err = mqttClient.Connect()
	if err != nil {
		log.Error(err)
	} else {
		inst.mqttClient = mqttClient
		inst.mqttConnected = true
	}

}

func (inst *Broker) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}

func (inst *Broker) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		go inst.setConnection()
	}
	if loopCount == 2 {
		go inst.subscribe()
	}
	if inst.mqttConnected {
		inst.WritePinTrue(node.Connected)
	} else {
		inst.WritePinFalse(node.Connected)
	}
	if loopCount > 1 {
		go inst.publish()
	}
}
