package broker

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const brokerHelp = `
A node used for pub/sub to an MQTT broker.
Please add a connection and select the connection to the broker`

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
	body = node.Defaults(body, mqttBroker, category)
	inputs := node.BuildInputs()
	outputs := node.BuildOutputs(node.BuildOutput(node.Connected, node.TypeBool, nil, body.Outputs, node.SetOutputHelp(node.ConnectedHelp)))
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(brokerHelp)
	network := &Broker{body, false, 0, nil, nil, false}
	network.Spec = body
	return network, nil
}

func (inst *Broker) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		go inst.setConnection()
	}
	if loopCount == 4 {
		go inst.subscribe()
	}
	if inst.mqttConnected {
		inst.WritePinTrue(node.Connected)
	} else {
		inst.WritePinFalse(node.Connected)
	}
	if loopCount > 5 {
		go inst.publish()
	}
}

func (inst *Broker) subscribe() {
	s := inst.GetStore()
	callback := func(client mqtt.Client, message mqtt.Message) {
		children, ok := s.Get(inst.GetID())
		payloads := getPayloads(children, ok)
		for _, payload := range payloads {
			if payload.Topic == message.Topic() {
				n := inst.GetNode(payload.NodeUUID)
				n.SetPayload(&node.Payload{
					Any: string(message.Payload()),
				})
			}
		}
	}
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	for _, payload := range payloads {
		if payload.IsPublisher {
			continue
		}
		if payload.Topic != "" {
			if inst.mqttClient != nil {
				err := inst.mqttClient.Subscribe(payload.Topic, mqttQOS, callback)
				if err != nil {
					log.Errorf("mqtt-broker subscribe:%s err:%s", payload.Topic, err.Error())
				} else {
					log.Infof("mqtt-broker subscribe:%s", payload.Topic)
				}
			}
		}
	}
}

type parsedBody struct {
	Body  interface{} `json:"body"`
	Topic string      `json:"topic"`
}

func parseBodyString(bodyString string) (*parsedBody, error) {
	body := &parsedBody{}
	err := json.Unmarshal([]byte(bodyString), &body)
	return body, err
}

func json2Str(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	return string(b), nil
}

func (inst *Broker) publish() {
	s := inst.GetStore()
	children, ok := s.Get(inst.GetID())
	payloads := getPayloads(children, ok)
	for _, payload := range payloads {
		if !payload.IsPublisher {
			continue
		}
		var bodyPayload string
		var useBodyTopic bool
		var mqttTopic string
		n := inst.GetNode(payload.NodeUUID)
		mqttTopic, _ = n.ReadPinAsString(node.Topic)
		p, _ := n.ReadPinAsString(node.Message)

		if mqttTopic == "" {
			body, err := parseBodyString(p)
			if err != nil {
				log.Errorf("mqtt-broker parse body err:%s", err.Error())
				continue
			}
			if body == nil {
				log.Errorf("mqtt-broker parse body is empty")
				continue
			}
			mqttTopic = body.Topic // use topic from payload
			useBodyTopic = true
			if mqttTopic == "" {
				log.Errorf("mqtt-broker topic can not be empty")
				continue
			}
			bodyAsString, err := json2Str(body.Body)
			if body == nil {
				log.Errorf("mqtt-broker convert body to string err:%s", err.Error())
				continue
			}
			bodyPayload = bodyAsString
		} else {
			bodyPayload = p
		}
		if mqttTopic == "" {
			log.Errorf("mqtt-broker topic can not be empty")
			continue
		}
		if p != "" {
			if mqttTopic == payload.Topic || useBodyTopic {
				if inst.mqttClient != nil {
					updated, _ := n.InputUpdated(node.Message)
					if updated {
						err := inst.mqttClient.Publish(mqttTopic, mqttQOS, mqttRetain, bodyPayload)
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

func getPayloads(children interface{}, ok bool) []*mqttPayload {
	if ok {
		mqttData := children.(*mqttStore)
		if mqttData != nil {
			return mqttData.payloads
		}
	}
	return nil
}
