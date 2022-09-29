package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	edge28lib "github.com/NubeDev/flow-eng/services/edge28"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	rubixIO "github.com/NubeDev/flow-eng/services/rubixio"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	*node.Spec
	client        *mqttclient.Client
	clients       *clients
	firstLoop     bool
	pingFailed    bool
	reconnectedOk bool
}

type clients struct {
	rio    *rubixIO.RubixIO
	edge28 *edge28lib.Edge28
}

func buildSubNodes(body *node.Spec, childNodes []*node.Spec) *node.Spec {
	for _, childNode := range childNodes {
		pprint.PrintJOSN(childNode)
	}
	body.SubFlow.Nodes = childNodes
	return body
}

var db *points.Store
var client *mqttclient.Client
var inst *Server
var application = names.Edge

//func NewServer(body *node.Spec, childNodes ...*node.Spec) (node.Node, error) {

func NewServer(body *node.Spec, store *points.Store) (node.Node, error) {
	var err error
	body = node.Defaults(body, server, category)
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
	client, err = mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err = client.Connect()
	if err != nil {
		log.Error(err)
		//return nil, err
	}
	c := &clients{}
	s := &Server{body, client, c, false, false, false}

	if application == names.RubixIO || application == names.RubixIOAndModbus {
		rubixIOUICount, rubixIOUOCount := points.CalcPointCount(1, application)
		rio := &rubixIO.RubixIO{}
		rio = rubixIO.New(&rubixIO.RubixIO{
			IP:          "0.0.0.0",
			StartAddrUI: rubixIOUICount,
			StartAddrUO: rubixIOUOCount,
			StartAddrDO: 2,
		})
		s.clients.rio = rio
	}
	if application == names.Edge {
		edge28 := edge28lib.New("192.168.15.141")
		s.clients.edge28 = edge28
	}
	db = store
	inst = s
	return s, err
}

func getServer() *Server {
	return inst
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
	err = client.Ping()
	if err != nil {
		log.Errorf("bacnet-server mqtt ping failed")
		inst.pingFailed = true
		err = getMqtt().Connect()
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

func getMqtt() *mqttclient.Client {
	return client
}

func getStore() *points.Store {
	if db == nil {
		log.Error("bacnet-server-node: store can not be empty")
		db = points.New(application, nil, 1, 200, 200)
	}
	return db
}
func getApplication() names.ApplicationName {
	return db.GetApplication()
}

func (inst *Server) subscribeBroker(topic string) {
	err := getMqtt().Subscribe(topic, mqttclient.AtLeastOnce, bacnetBus)
	if err != nil {
		log.Errorf("bacnet-server mqtt:%s", err.Error())
		inst.pingFailed = false
	}
}

func (inst *Server) subscribeToRubixIO() {
	if getApplication() == names.RubixIO {
		err := getMqtt().Subscribe("rubixio/inputs/all", mqttclient.AtLeastOnce, rubixIOBus)
		if err != nil {
			log.Errorf("bacnet-server mqtt:%s", err.Error())
			inst.pingFailed = false
		}
	}
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
