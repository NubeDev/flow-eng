package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

type NodeA struct {
	*node.BaseNode
}

func Builder(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case "nodeA":
		return NewNodeA(body)
	case "nodeB":
		//return NewNodeB(body)
	}
	return nil, errors.New("node not found")
}

func NewNodeA(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body)
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, body.Inputs), node.BuildInput(node.In2, node.TypeFloat, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, body.Outputs))
	return &NodeA{body}, nil
}

func (n *NodeA) Process() {
	_, r := n.ReadPin(node.In1)
	fmt.Println("READ IN-1", n.GetNodeName(), r)

	_, in1Val, in1Not := n.ReadPinNum(node.In1)
	_, in2Val, in2Not := n.ReadPinNum(node.In2)

	if in1Not && in2Not {
		add := in1Val + in2Val
		fmt.Println(add, "WRITE----------", n.GetNodeName(), in1Val+in2Val)
		n.WritePinNum(node.Out1, add)
		return
	}
	if in1Not {
		n.WritePinNum(node.Out1, in1Val)
		return
	}
	if in2Not {
		n.WritePinNum(node.Out1, in2Val)
		return
	}
}

func (n *NodeA) Cleanup() {}
