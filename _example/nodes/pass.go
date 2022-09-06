package nodes

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/node"
)

// a node out should be able to have multiple connections
// if we ref a Connection by the name it's easier to manage migrations but will be more work in coding the app

type Connection struct {
	NodeID string `json:"nodeID"`
	Port   string `json:"port"`
}

type PortCommon struct {
	Name       portName    `json:"name"` // in1
	Type       dataTypes   `json:"type"` // int8
	Connection *Connection `json:"connection"`
}

type PortCommonOut struct {
	Name        portName      `json:"name"` // in1
	Type        dataTypes     `json:"type"` // int8
	Connections []*Connection `json:"connection"`
}

type TypeInput struct {
	*PortCommon
	*node.InputPort
	value *adapter.Int8
}

type TypeOutput struct {
	*PortCommonOut
	*node.OutputPort
	value *adapter.Int8
}

type dataTypes string
type portName string

const (
	typeInt8 dataTypes = "int8"
)

const (
	in1  portName = "in1"
	out1 portName = "out1"
)

type Node struct {
	NodeID  string        `json:"nodeID"` // abc
	Name    string        `json:"name"`   // my node
	Node    string        `json:"node"`   // PASS
	Inputs  []*TypeInput  `json:"inputs"`
	Outputs []*TypeOutput `json:"outputs"`
	info    node.NodeInfo
}

func buildInput(inputType dataTypes, conn *Connection) *TypeInput {
	var dataType buffer.Type
	if inputType == typeInt8 {
		dataType = buffer.Int8
	}
	var port = node.NewInputPort(dataType, nil)
	return &TypeInput{
		PortCommon: &PortCommon{
			Name: in1,
			Type: typeInt8,
			Connection: &Connection{
				NodeID: conn.NodeID,
				Port:   conn.Port,
			},
		},
		InputPort: port,
		value:     adapter.NewInt8(port),
	}
}

func buildOutput(inputType dataTypes, conn []*Connection) *TypeOutput {
	var dataType buffer.Type
	if inputType == typeInt8 {
		dataType = buffer.Int8
	}
	var port = node.NewOutputPort(dataType, nil)
	return &TypeOutput{
		PortCommonOut: &PortCommonOut{
			Name:        out1,
			Type:        typeInt8,
			Connections: conn,
		},
		OutputPort: port,
		value:      adapter.NewInt8(port),
	}
}

func getInput(body *Node, num int) *TypeInput {
	for i, input := range body.Inputs {
		if i == num {
			return input
		}
	}
	return nil
}

func New(body *Node) *Node {
	buildIn1 := buildInput(typeInt8, getInput(body, 0).Connection)
	buildOut1 := buildOutput(typeInt8, body.Outputs[0].Connections)
	return &Node{
		NodeID:  "1111",
		Name:    body.Name,
		Node:    "PASS",
		Inputs:  []*TypeInput{buildIn1},
		Outputs: []*TypeOutput{buildOut1},
	}
}

func (node *Node) Get() node.NodeInfo {
	return node.info
}

func (node *Node) Info() node.NodeInfo {
	return node.info
}

func (node *Node) Process() {

	//read := node.reader.Get()
	//if node.Info().Name == "nodeA" {
	//	fmt.Println("get reader", node.Info().Name, node.reader.Get())
	//	node.writer.Set(11)
	//	fmt.Println("get writer value", node.Info().Name, node.writer.Get())
	//} else {
	//	fmt.Println(node.Info().Name, read)
	//}

}

func (node *Node) Cleanup() {}
