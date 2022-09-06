package node

import "fmt"

type NodeB struct {
	*Node
}

func (n *NodeB) Process() {
	fmt.Println("Node B>>>", n.GetName())
}

func (n *NodeB) Cleanup() {

}
