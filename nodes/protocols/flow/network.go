package flow

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type Network struct {
	*node.Spec
	firstLoop     bool
	loopCount     uint64
	connection    *db.Connection
	mqttClient    *mqttclient.Client
	mqttConnected bool
	points        []*point
}

var mqttQOS = mqttclient.AtMostOnce
var mqttRetain = false

func NewNetwork(body *node.Spec) (node.Node, error) {
	//var err error
	body = node.Defaults(body, flowNetwork, category)
	inputs := node.BuildInputs()
	outputs := node.BuildOutputs(node.BuildOutput(node.Connected, node.TypeBool, nil, body.Outputs))
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	network := &Network{body, false, 0, nil, nil, false, nil}
	return network, nil
}

func getPayloads(children interface{}, ok bool) []*pointDetails {
	if ok {
		mqttData := children.(*pointStore)
		if mqttData != nil {
			return mqttData.payloads
		}
	}
	return nil
}

func (inst *Network) subscribe() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		s := inst.GetStore()
		children, ok := s.Get(inst.GetID())
		payloads := getPayloads(children, ok)
		for _, payload := range payloads {
			if payload.topic == message.Topic() {
				n := inst.GetNode(payload.nodeUUID)
				n.SetPayload(&node.Payload{
					String: str.New(string(message.Payload())),
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
					log.Errorf("flow-network-broker subscribe:%s err:%s", payload.topic, err.Error())
				} else {
					log.Infof("flow-network-broker subscribe:%s", payload.topic)
				}
			}
		}
	}
}

func (inst *Network) pointsList() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		var points []*point
		err := json.Unmarshal(message.Payload(), &points)
		if err == nil {
			if points != nil {
				s := inst.GetStore()
				s.Set(fmt.Sprintf("pointsList_%s", inst.GetID()), points, 0)
			}
		} else {
			log.Errorf("failed to get flow-framework points list err:%s", err.Error())
		}
	}
	var topic = "+/+/+/+/+/+/rubix/points/value/points"
	if inst.mqttClient != nil {
		err := inst.mqttClient.Subscribe(topic, mqttQOS, callback)
		if err != nil {
			log.Errorf("flow-network-broker subscribe:%s err:%s", topic, err.Error())
		} else {
			log.Infof("flow-network-broker subscribe:%s", topic)
		}
	}
}

func (inst *Network) publish() {
	s := inst.GetStore()
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	for _, payload := range payloads {
		if !payload.isWriteable {
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
							log.Errorf("flow-network-broker publish err:%s", err.Error())
						}
					}

				}
			}
		}
	}
}

func (inst *Network) setConnection() {
	settings, err := getSettings(inst.GetSettings())
	if err != nil {
		log.Errorf("add mqtt broker failed to get settings err:%s", err.Error())
		return
	}
	connection, err := inst.GetDB().GetConnection(settings.Conn)
	if err != nil {
		log.Error("flow-network add mqtt broker failed to find connection")
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

func (inst *Network) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}

func (inst *Network) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		go inst.setConnection()
	}
	if loopCount == 2 {
		go inst.subscribe()
		go inst.pointsList()
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
