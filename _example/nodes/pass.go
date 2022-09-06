package nodes

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

// a node out should be able to have multiple connections
// if we ref a Connection by the name it's easier to manage migrations but will be more work in coding the app

type Node struct {
	InputList  []*node.TypeInput  `json:"inputs"`
	OutputList []*node.TypeOutput `json:"outputs"`
	Info       node.Info          `json:"info"`
}

const (
	nodeType    = "pass"
	inputCount  = 1
	outputCount = 1
)

func New(body *Node) (*Node, error) {
	body, err := Check(body, NodeSpec{nodeType, inputCount, outputCount})
	if err != nil {
		return nil, err
	}

	return &Node{
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

func (n *Node) GetName() string {
	return n.GetInfo().Name
}

func (n *Node) GetID() string {
	return n.GetInfo().NodeID
}

func (n *Node) GetType() string {
	return n.GetInfo().Type
}

func (n *Node) GetInputs() []*node.TypeInput {
	return n.InputList
}

func (n *Node) GetOutputs() []*node.TypeOutput {
	return n.OutputList
}

func (n *Node) GetInfo() node.Info {
	return n.Info
}

func (n *Node) Process() {

	for _, input := range n.GetInputs() {
		fmt.Println("READ-VALUE", input.ValueFloat64.Get(), "NAME", n.Info.Name)
	}

	for _, output := range n.GetOutputs() {
		if n.Info.Name == "nodeA" {
			fmt.Println("WRITE-VALUE", output.ValueFloat64.Get(), "NAME", n.Info.Name)
			output.ValueFloat64.Set(33.3)

		}

	}
}

func (n *Node) Cleanup() {}
