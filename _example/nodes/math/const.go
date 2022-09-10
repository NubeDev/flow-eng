package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Const struct {
	*node.BaseNode
}

func NewConst(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body, constNum)
	body.Info.Name = constNum
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, 0, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, body.Outputs))
	return &Const{body}, nil
}

func (inst *Const) Process() {
	_, in1Val, in1Not := inst.ReadPinNum(node.In1)
	if in1Not {
		inst.WritePinNum(node.Out1, in1Val)
	}

}

func (inst *Const) Cleanup() {}
