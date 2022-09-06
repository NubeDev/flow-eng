package back2

import (
	"fmt"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/node"
)

type Pass struct {
	info node.NodeInfo

	Input  *node.InputPort
	Output *node.OutputPort

	reader *adapter.Int8
	writer *adapter.Int8
}

func (node *Pass) AddConnection(inputs ...*node.InputPort) {
	//TODO implement me
	//panic("implement me")
}

func NewPass(name string) *Pass {
	input := node.NewInputPort(buffer.Int8, nil)
	output := node.NewOutputPort(buffer.Int8, nil)
	info := node.NodeInfo{Name: name, Description: "passes input dat to another node", Version: "1.0.0"}
	return &Pass{info, input, output, adapter.NewInt8(input), adapter.NewInt8(output)}
}

func (node *Pass) Info() node.NodeInfo {
	return node.info
}

func (node *Pass) Process() {
	read := node.reader.Get()
	if node.Info().Name == "nodeA" {
		fmt.Println("get reader", node.Info().Name, node.reader.Get())
		node.writer.Set(11)
		fmt.Println("get writer value", node.Info().Name, node.writer.Get())
	} else {
		fmt.Println(node.Info().Name, read)
	}

}

func (node *Pass) Cleanup() {}
