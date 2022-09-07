package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
)

type NodeA struct {
	*Node
}

func SpecNodeA(body *Node) (*Node, error) {
	body = emptyNode(body)
	body.Info.Name = "nodeA"
	body.Info.NodeID = setUUID(body.Info.NodeID)
	body.Inputs = buildInputs(buildInput(In1, TypeFloat64, body.Inputs), buildInput(In2, TypeFloat64, body.Inputs))
	body.Outputs = buildOutputs(buildOutput(Out1, TypeFloat64, body.Outputs))
	return body, nil
}

func (n *NodeA) Process() {

	_, r := n.readPinValue(In1)
	fmt.Println("READ IN-1", n.GetNodeName(), r)

	in1, in1Val := n.readPinValue(In1)
	in2, in2Val := n.readPinValue(In2)

	if !float.IsNil(in1) && !float.IsNil(in2) {
		val := float.New(in1Val + in2Val)
		fmt.Println(val, "WRITE----------", n.GetNodeName(), in1Val+in2Val)
		n.writePinValue(Out1, val)
		return
	}
	if !float.IsNil(in1) {
		val := float.New(in1Val)
		n.writePinValue(Out1, val)
		return
	}
	if !float.IsNil(in2) {
		val := float.New(in2Val)
		n.writePinValue(Out1, val)
		return
	}

}

func (n *NodeA) Cleanup() {}
