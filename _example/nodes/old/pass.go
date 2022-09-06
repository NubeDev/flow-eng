package nodes

import (
	"fmt"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/node"
)

type typeInput struct {
	*node.InputPort
	reader *adapter.Int8
}

type typeOutput struct {
	*node.OutputPort
	writer *adapter.Int8
}

type Pass struct {
	node.NodeInfo
	Input  *typeInput  `json:"input"`
	Output *typeOutput `json:"output"`
}

type nodeBody struct {
}

func NewPass(name string) *Pass {
	var inputPort = node.NewInputPort(buffer.Int8, &node.Details{
		Name:               "in1",
		DataType:           "int8",
		ConnectionNodeName: "",
		ConnectionName:     "",
	})
	input := &typeInput{
		InputPort: inputPort,
		reader:    adapter.NewInt8(inputPort),
	}

	outputPort := node.NewOutputPort(buffer.Int8, &node.Details{
		Name:               "out1",
		DataType:           "int8",
		ConnectionNodeName: "",
		ConnectionName:     "",
	})
	output := &typeOutput{
		OutputPort: outputPort,
		writer:     adapter.NewInt8(outputPort),
	}

	info := node.NodeInfo{Name: name, Type: "PASS", Description: "passes input dat to another node", Version: "1.0.0"}
	return &Pass{info, input, output}
}

func (node *Pass) Info() node.NodeInfo {
	return node.NodeInfo
}

func (node *Pass) Process() {

	if node.Name == "nodeA" {
		fmt.Println("Process", node.Name, node.Input.reader.Get())
		node.Output.writer.Set(22)
		fmt.Println("Get writer value", node.Name, node.Output.writer.Get(), node.Output.UUID())
	} else {

		fmt.Println("Process", node.Name, node.Input.reader.Get(), node.Input.UUID())

	}

}

func (node *Pass) Cleanup() {}
