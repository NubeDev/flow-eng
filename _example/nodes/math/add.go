package math

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

type Add struct {
	*node.BaseNode
}

func NewAdd(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body)
	body.Info.Name = add
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, body.Inputs), node.BuildInput(node.In2, node.TypeFloat, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, body.Outputs))
	return &Add{body}, nil
}

func (inst *Add) Process() {
	_, r := inst.ReadPin(node.In1)
	fmt.Println("READ IN-1", inst.GetNodeName(), r)
	_, in1Val, in1Not := inst.ReadPinNum(node.In1)
	_, in2Val, in2Not := inst.ReadPinNum(node.In2)

	if in1Not && in2Not {
		add := in1Val + in2Val
		fmt.Println(add, "WRITE----------", inst.GetNodeName(), in1Val+in2Val)
		inst.WritePinNum(node.Out1, add)
		return
	}
	if in1Not {
		inst.WritePinNum(node.Out1, in1Val)
		return
	}
	if in2Not {
		inst.WritePinNum(node.Out1, in2Val)
		return
	}
}

func (inst *Add) Cleanup() {}
