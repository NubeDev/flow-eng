package bacnet

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

type Server struct {
	*node.Spec
}

var db *bstore.BacnetStore

func buildSubNodes(body *node.Spec, childNodes []*node.Spec) *node.Spec {
	for _, childNode := range childNodes {
		pprint.PrintJOSN(childNode)
	}
	body.SubFlow.Nodes = childNodes
	return body
}

func NewServer(body *node.Spec, childNodes ...*node.Spec) (node.Node, error) {
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
	body = node.BuildNode(body, nil, outputs, nil)

	application := applications.Edge // make this a setting eg: if it's an edge-28 it would give the user 8AI, 8AOs and 100 BVs/AVs
	//if application == "" {
	//	application = applications.BACnet
	//}
	db = bstore.New(application, nil)
	return &Server{body}, nil
}

func GetStore() *bstore.BacnetStore {
	return db
}

func (inst *Server) GetStore() *bstore.BacnetStore {
	return inst.db()
}

func (inst *Server) db() *bstore.BacnetStore {
	return db
}

func (inst *Server) Process() {

}

func (inst *Server) Cleanup() {}
