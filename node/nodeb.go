package node

import "fmt"

type NodeB struct {
	*BaseNode
}

func NewNodeB(body *BaseNode) (Node, error) {
	body = emptyNode(body)
	body.Info.NodeID = setUUID(body.Info.NodeID)
	body.Inputs = buildInputs(buildInput("in1", TypeFloat64, body.Inputs), buildInput("in2", TypeFloat64, body.Inputs))
	body.Outputs = buildOutputs(buildOutput("out1", TypeFloat64, body.Outputs))
	return &NodeB{BaseNode: body}, nil
}

func (n *NodeB) Process() {
	for _, input := range n.GetInputs() {
		fmt.Println("Node B>>>", n.GetName(), input.ValueFloat64.Get())
	}
}

func (n *NodeB) Cleanup() {

}
