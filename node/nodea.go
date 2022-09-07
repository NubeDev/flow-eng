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
	body.Inputs = buildInputs(buildInput("in1", TypeFloat64, body.Inputs), buildInput("in2", TypeFloat64, body.Inputs))
	body.Outputs = buildOutputs(buildOutput("out1", TypeFloat64, body.Outputs))
	return body, nil
}

func (n *NodeA) Process() {

	v := n.readPinValue("in1")
	fmt.Println("read", float.NonNil(v), "NAME", n.GetNodeName())
	if n.GetNodeName() == "name-a123" {
		n.writePinValue("out1", 222)
	}

}

func (n *NodeA) Cleanup() {}
