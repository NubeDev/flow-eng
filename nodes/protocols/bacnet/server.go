package bacnet

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	eventbus "github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	rubixIO "github.com/NubeDev/flow-eng/services/rubixio"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	*node.Spec
	client        *mqttclient.Client
	rio           *rubixIO.RubixIO
	firstLoop     bool
	pingFailed    bool
	reconnectedOk bool
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

func NewServer(body *node.Spec, childNodes ...*node.Spec) (node.Node, error) {
	var err error
	body = node.Defaults(body, server, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputBroker := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	outputApplication := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	outputErr := node.BuildOutput(node.ErrMsg, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(outputBroker, outputApplication, outputErr)
	parameters := &node.Parameters{
		Application: &node.Application{
			Application: applications.BACnet,
			IsChild:     false,
		},
		MaxNodeCount: 1,
	}
	body.Parameters = node.BuildParameters(parameters) // if node is already added then show the user
	body = buildSubNodes(body, childNodes)
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, nil)
	application := applications.RubixIO // make this a setting eg: if it's an edge-28 it would give the user 8AI, 8AOs and 100 BVs/AVs
	client, err = mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err = client.Connect()
	if err != nil {
		log.Error(err)
		//return nil, err
	}
	eventbus.New()
	rio := rubixIO.New()
	db = points.New(application, nil)
	s := &Server{body, client, rio, false, false, false}
	inst = s
	return s, err
}

func getServer() *Server {
	return inst
}

func (inst *Server) intBus() {
	go inst.priorityBus()
	go inst.rubixIOBus()
}

func (inst *Server) Process() {
	inst.intBus()

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
		panic("bacnet-server-node: store can not be empty")
	}
	return db
}
func getApplication() node.ApplicationName {
	return db.GetApplication()
}

func (inst *Server) subscribeBroker(topic string) {
	err := getMqtt().Subscribe(topic, mqttclient.AtLeastOnce, eventbus.PointsHandler)
	if err != nil {
		log.Errorf("bacnet-server mqtt:%s", err.Error())
		inst.pingFailed = false
	}
}

func (inst *Server) subscribeToRubixIO() {
	if getApplication() == applications.RubixIO {
		inst.subscribeBroker("rubixio/inputs/all")
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
