package math

import (
	"fmt"
	"github.com/NubeDev/flow-eng/_example/nodes"
	"github.com/NubeDev/flow-eng/node"
)

type Node struct {
	InputList  []*node.TypeInput  `json:"inputs"`
	OutputList []*node.TypeOutput `json:"outputs"`
	Info       node.Info          `json:"info"`
}

const (
	category    = "math"
	nodeType    = "add"
	inputCount  = 2
	outputCount = 2
)

func New(body *nodes.Node) (*Node, error) {
	//body, err := nodes.Check(body, nodes.NodeSpec{NodeType: nodeType, InputCount: inputCount, OutputCount: outputCount})
	//if err != nil {
	//	return nil, err
	//}
	//
	//buildIn1 := nodes.BuildInput(node.TypeInt8, nodes.GetInput(body, 0).Connection)
	//buildOut1 := nodes.BuildOutput(node.TypeInt8, nodes.GetOutConnections(body, 0))
	//
	//return &Node{
	//	InputList:  []*node.TypeInput{buildIn1},
	//	OutputList: []*node.TypeOutput{buildOut1},
	//	Info: node.NodeInfo{
	//		NodeID:      body.Info.NodeID,
	//		Name:        body.Info.Name,
	//		Type:        nodeType,
	//		Category:    category,
	//		Description: "desc",
	//		Version:     "1",
	//	},
	//}, nil

	return nil, nil
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
