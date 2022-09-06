package node

import "fmt"

type NodeA struct {
	*Node
}

func (n *NodeA) Process() {
	fmt.Println("Node A>>>", n.GetName())
}

func (n *NodeA) Cleanup() {

}
