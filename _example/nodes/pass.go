package nodes

import (
	"fmt"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/node"
)

// a node out should be able to have multiple connections
// if we ref a Connection by the name it's easier to manage migrations but will be more work in coding the app

type Node struct {
	NodeID     string             `json:"nodeID"` // abc
	Name       string             `json:"name"`   // my node
	NodeType   string             `json:"node"`   // PASS
	InputList  []*node.TypeInput  `json:"inputs"`
	OutputList []*node.TypeOutput `json:"outputs"`
	info       node.NodeInfo
}

func (n *Node) Inputs() []*node.TypeInput {
	return n.InputList
}

func (n *Node) Outputs() []*node.TypeOutput {
	return n.OutputList
}

func buildInput(inputType node.DataTypes, conn *node.Connection) *node.TypeInput {
	var dataType buffer.Type
	if inputType == node.TypeInt8 {
		dataType = buffer.Int8
	}
	var port = node.NewInputPort(dataType, nil)
	return &node.TypeInput{
		PortCommon: &node.PortCommon{
			Name: node.In1,
			Type: node.TypeInt8,
			Connection: &node.Connection{
				NodeID: conn.NodeID,
				Port:   conn.Port,
			},
		},
		InputPort: port,
		Value:     adapter.NewInt8(port),
	}
}

func buildOutput(inputType node.DataTypes, conn []*node.Connection) *node.TypeOutput {
	var dataType buffer.Type
	if inputType == node.TypeInt8 {
		dataType = buffer.Int8
	}
	var port = node.NewOutputPort(dataType, nil)
	return &node.TypeOutput{
		PortCommonOut: &node.PortCommonOut{
			Name:        node.Out1,
			Type:        node.TypeInt8,
			Connections: conn,
		},
		OutputPort: port,
		Value:      adapter.NewInt8(port),
	}
}

func getInput(body *Node, num int) *node.TypeInput {
	for i, input := range body.InputList {
		if i == num {
			return input
		}
	}
	return nil
}

func getOutConnections(body *Node, num int) []*node.Connection {
	for _, output := range body.OutputList {
		return output.Connections
	}
	return nil
}

func New(body *Node) *Node {
	buildIn1 := buildInput(node.TypeInt8, getInput(body, 0).Connection)
	buildOut1 := buildOutput(node.TypeInt8, getOutConnections(body, 0))

	return &Node{
		NodeID:     body.NodeID,
		Name:       body.Name,
		NodeType:   "PASS",
		InputList:  []*node.TypeInput{buildIn1},
		OutputList: []*node.TypeOutput{buildOut1},
		info: node.NodeInfo{
			Name:        body.Name,
			Type:        "PASS",
			Description: "desc",
			Version:     "1",
		},
	}
}

func (n *Node) Info() node.NodeInfo {
	return n.info
}

func (n *Node) Process() {

	for _, input := range n.InputList {
		fmt.Println("READ-VALUE", input.Value.Get(), "NAME", n.info.Name)
	}

	for _, output := range n.OutputList {
		if n.info.Name == "a123" {
			fmt.Println("WRITE-VALUE", output.Value.Get(), "NAME", n.info.Name)
			output.Value.Set(11)

		}

	}

	//read := node.reader.Get()
	//if node.Info().Name == "nodeA" {
	//	fmt.Println("get reader", node.Info().Name, node.reader.Get())
	//	node.writer.Set(11)
	//	fmt.Println("get writer value", node.Info().Name, node.writer.Get())
	//} else {
	//	fmt.Println(node.Info().Name, read)
	//}

}

func (n *Node) Cleanup() {}
