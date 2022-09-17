package bacnet

import (
	"errors"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/helpers/store"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
)

type Server struct {
	*node.Spec
}

var db *store.Store

func NewServer(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, server, category)
	outputs := node.BuildOutputs(node.BuildOutput(node.ErrMsg, node.TypeString, nil, body.Outputs))
	body.Parameters = node.BuildParameters(&node.Parameters{Application: applications.BACnet, MaxNodeCount: 1}) // if node is already added then show the user
	body = node.BuildNode(body, nil, outputs, nil)
	db = store.Init()
	return &Server{body}, nil
}

func DB() (*store.Store, error) {
	if db == nil {
		return nil, errors.New("bacnet store has not be initialised")
	}
	return db, nil
}

func (inst *Server) db() *store.Store {
	return db
}

func (inst *Server) Process() {

	inst.db().Set("bacnet", &Point{ObjectType: "aaaaaaaa"}, 0)

	cli, ok := db.Get("bacnet")
	if !ok {

	}
	parse := cli.(*Point)
	pprint.Print(parse)
}

func (inst *Server) Cleanup() {}
