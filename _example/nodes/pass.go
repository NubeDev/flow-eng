package nodes

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

// a node out should be able to have multiple connections
// if we ref a Connection by the name it's easier to manage migrations but will be more work in coding the app

type Node struct {
	Node *node.Node
}

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

func (n *Node) Process() {

	n.Node.Process()

	fmt.Println("WANT here", n.Node.Info.Name)
	fmt.Println(n.Node.GetName(), 9999999999)
	//for _, input := range n.GetInputs() {
	//	fmt.Println("READ-VALUE", input.ValueFloat64.Get(), "NAME", n.Info.Name)
	//}
	//
	//for _, output := range n.GetOutputs() {
	//	if n.Info.Name == "nodeA" {
	//		fmt.Println("WRITE-VALUE", output.ValueFloat64.Get(), "NAME", n.Info.Name)
	//		output.ValueFloat64.Set(33.3)
	//
	//	}
	//
	//}
}

func (n *Node) Cleanup() {}
