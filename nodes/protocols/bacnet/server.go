package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	edge28lib "github.com/NubeDev/flow-eng/services/edge28"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	rubixIO "github.com/NubeDev/flow-eng/services/rubixio"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type Bacnet struct {
	Store       *points.Store
	MqttClient  *mqttclient.Client
	Application names.ApplicationName
}

type Server struct {
	*node.Spec
	clients       *clients
	firstLoop     bool
	pingFailed    bool
	reconnectedOk bool
	store         *points.Store
	application   names.ApplicationName
}

type clients struct {
	mqttClient *mqttclient.Client
	rio        *rubixIO.RubixIO
	edge28     *edge28lib.Edge28
}

//var db *points.Store
//var mqttClient *mqttclient.Client
//var inst *Server
//var application = names.Edge

func NewServer(body *node.Spec, opts *Bacnet) (node.Node, error) {
	var application = opts.Application
	application = names.Edge
	var err error
	ip := "192.168.15.141"
	body = node.Defaults(body, serverNode, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputBroker := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	outputApplication := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	outputErr := node.BuildOutput(node.ErrMsg, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(outputBroker, outputApplication, outputErr)
	parameters := &node.Parameters{
		Application: &node.Application{
			Application: names.BACnet,
			IsChild:     false,
		},
		MaxNodeCount: 1,
	}
	body.Parameters = node.BuildParameters(parameters) // if node is already added then show the user
	//body = buildSubNodes(body, childNodes)
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, nil)

	clients := &clients{}
	server := &Server{body, clients, false, false, false, opts.Store, opts.Application}
	server.clients.mqttClient = opts.MqttClient
	if application == names.RubixIO || application == names.RubixIOAndModbus {
		rubixIOUICount, rubixIOUOCount := points.CalcPointCount(1, application)
		rio := &rubixIO.RubixIO{}
		rio = rubixIO.New(&rubixIO.RubixIO{
			IP:          ip,
			StartAddrUI: rubixIOUICount,
			StartAddrUO: rubixIOUOCount,
			StartAddrDO: 2,
		})
		server.clients.rio = rio
		log.Infof("bacnet-server: start application: %s device-ip: %s", application, ip)
	}
	if application == names.Edge {
		edge28 := edge28lib.New(ip)
		server.clients.edge28 = edge28
		log.Infof("bacnet-server: start application: %s device-ip: %s", application, ip)
	}

	return server, err
}

func (inst *Server) Process() {
	if !inst.firstLoop {
		go inst.subscribeToRubixIO()
		inst.firstLoop = true
	}
	if inst.pingFailed || inst.reconnectedOk { // on failed resubscribe
		go inst.subscribeToRubixIO()
	}
	go inst.mqttReconnect()
	go inst.protocolRunner()

}

func (inst *Server) mqttReconnect() {
	var err error
	err = inst.clients.mqttClient.Ping()
	if err != nil {
		log.Errorf("bacnet-server mqtt ping failed")
		inst.pingFailed = true
		err = inst.clients.mqttClient.Connect()
		if err != nil {
			log.Errorf("bacnet-server failed to reconnect with mqtt broker")
			inst.reconnectedOk = false
		} else {
			inst.reconnectedOk = true
		}
	} else {
		fmt.Println("ping server ok")
	}
}

//func getMqtt() *mqttclient.Client {
//	return mqttClient
//}
//
//func getStore() *points.Store {
//	if db == nil {
//		//log.Error("bacnet-server-node: store can not be empty")
//		db = points.New(application, nil, 1, 200, 200)
//	}
//	return db
//}
//func getApplication() names.ApplicationName {
//	if db != nil {
//		return db.GetApplication()
//	}
//	return ""
//}

func (inst *Server) subscribeBroker(topic string) {
	callback := func(client mqtt.Client, message mqtt.Message) {
		rawData := message.Payload()
		fmt.Println(message.Topic(), string(rawData))
	}
	err := inst.clients.mqttClient.Subscribe("test", mqttclient.AtLeastOnce, callback)
	if err != nil {
		log.Errorf("bacnet-server mqtt:%s", err.Error())
		inst.pingFailed = false
	}

}

func (inst *Server) subscribeToRubixIO() {
	//if inst.application == names.RubixIO {
	//	err := inst.clients.mqttClient.Subscribe("rubixio/inputs/all", mqttclient.AtLeastOnce, rubixIOBus)
	//	if err != nil {
	//		log.Errorf("bacnet-server mqtt:%s", err.Error())
	//		inst.pingFailed = false
	//	}
	//}
	objs := []string{"ai", "ao", "av", "bi", "bo", "bv"}
	for _, obj := range objs {
		topic := fmt.Sprintf("%s/+/pv", topicObjectBuilder(obj))
		inst.subscribeBroker(topic)
	}
	objsOuts := []string{"ao", "av", "bo", "bv"}
	for _, obj := range objsOuts {
		topic := fmt.Sprintf("%s/+/pri", topicObjectBuilder(obj))
		inst.subscribeBroker(topic)
	}
}

func (inst *Server) Cleanup() {}
