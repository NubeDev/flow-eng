package node

import "fmt"

type NodeA struct {
	*Node
}

func (n *NodeA) Process() {
	fmt.Println("Node A>>>", n.GetName())
	for _, input := range n.GetInputs() {
		fmt.Println(input.Name)
	}

	for _, out := range n.GetOutputs() {
		out.ValueFloat64.Set(22)
	}
}

func (n *NodeA) Cleanup() {}
