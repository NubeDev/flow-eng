package node

import "fmt"

type NodeB struct {
	*Node
}

func (n *NodeB) Process() {
	for _, input := range n.GetInputs() {
		fmt.Println("Node B>>>", n.GetName(), input.ValueFloat64.Get())
	}
}

func (n *NodeB) Cleanup() {

}
