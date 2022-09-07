package node

import (
	"fmt"
)

type NodeA struct {
	*BaseNode
}

func NewNodeA(body *BaseNode) (Node, error) {
	body = emptyNode(body)
	body.Info.NodeID = setUUID(body.Info.NodeID)
	body.Inputs = buildInputs(buildInput(In1, TypeFloat, body.Inputs), buildInput(In2, TypeFloat, body.Inputs))
	body.Outputs = buildOutputs(buildOutput(Out1, TypeFloat, body.Outputs))
	return &NodeA{BaseNode: body}, nil
}

func (n *NodeA) Process() {
	_, r := n.readPin(In1)
	fmt.Println("READ IN-1", n.GetNodeName(), r)

	_, in1Val, in1Not := n.readPinNum(In1)
	_, in2Val, in2Not := n.readPinNum(In2)

	if in1Not && in2Not {
		add := in1Val + in2Val
		fmt.Println(add, "WRITE----------", n.GetNodeName(), in1Val+in2Val)
		n.writePinNum(Out1, add)
		return
	}
	if in1Not {
		n.writePinNum(Out1, in1Val)
		return
	}
	if in2Not {
		n.writePinNum(Out1, in2Val)
		return
	}
}

func (n *NodeA) Cleanup() {}
