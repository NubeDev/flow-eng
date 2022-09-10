package math

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
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

func getPoints() {
	client := resty.New()
	resp, err := client.R().
		SetResult(&node.BaseNode{}).
		Get("http://192.168.15.190:1660/api/points")
	fmt.Println(err)
	fmt.Println(resp.Status())
	//fmt.Println(resp.String())
}

func (inst *Add) Process() {
	Process(inst)
	//go getPoints()

}

func (inst *Add) Cleanup() {}
