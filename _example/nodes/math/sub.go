package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Sub struct {
	*node.BaseNode
}

func NewSub(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body, sub)
	body.Info.Name = node.SetName(sub)
	body.Info.Category = node.SetName(category)
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.In2, node.TypeFloat, nil, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs))
	return &Sub{body}, nil
}

func (inst *Sub) Process() {
	Process(inst)
}

func (inst *Sub) Cleanup() {}
