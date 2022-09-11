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
	body = node.EmptyNode(body, add)
	body.Info.Name = node.SetName(add)
	body.Info.Category = node.SetName(category)
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.DynamicInputs(node.In, node.TypeFloat, nil, 2, 2, body.Inputs)...)
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs))
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
