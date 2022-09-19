package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/cbus"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

type Server struct {
	*node.Spec
	client *mqttbase.Mqtt
}

var db *bstore.BacnetStore

func buildSubNodes(body *node.Spec, childNodes []*node.Spec) *node.Spec {
	for _, childNode := range childNodes {
		pprint.PrintJOSN(childNode)
	}
	body.SubFlow.Nodes = childNodes
	return body
}

var client *mqttbase.Mqtt

func NewServer(body *node.Spec, childNodes ...*node.Spec) (node.Node, error) {
	var err error
	body = node.Defaults(body, server, category)
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
	body = node.BuildNode(body, nil, outputs, nil)

	application := applications.Modbus // make this a setting eg: if it's an edge-28 it would give the user 8AI, 8AOs and 100 BVs/AVs

	client, err = mqttbase.NewMqtt()
	client.Connect()

	db = bstore.New(application, nil)
	return &Server{body, client}, err
}

func getStore() *bstore.BacnetStore {
	if db == nil {
		panic("bacnet-server-node: store can not be empty")
	}
	return db
}
func getRunnerType() node.ApplicationName {
	return db.GetApplication()
}

func (inst *Server) db() *bstore.BacnetStore {
	return db
}

func (inst *Server) bus() cbus.Bus {
	return inst.client.BACnetBus()
}

func (inst *Server) Process() {
	inst.protocolRunner()
}

func (inst *Server) Cleanup() {}
