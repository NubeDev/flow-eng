package nodes

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

// a node out should be able to have multiple connections
// if we ref a Connection by the name it's easier to manage migrations but will be more work in coding the app

const (
	nodeType    = "pass"
	inputCount  = 1
	outputCount = 1
)

func New(body *node.Node) (*node.Node, error) {
	body, err := Check(body, NodeSpec{nodeType, inputCount, outputCount})
	if err != nil {
		return nil, err
	}
	fmt.Println(body.Info.Name, body.GetName())
	return &node.Node{
		InputList:  BuildInputs(body),
		OutputList: BuildOutputs(body),
		Info: node.Info{
			NodeID:      body.Info.NodeID,
			Name:        body.Info.Name,
			Type:        nodeType,
			Description: "desc",
			Version:     "1",
		},
	}, nil
}
