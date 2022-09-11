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
	body = node.Defaults(body, add, category)
	buildCount, setting, count, err := inputsCount(body)
	if err != nil {
		return nil, err
	}
	settings, err := node.BuildSettings(setting)
	if err != nil {
		return nil, err
	}
	inputs := node.BuildInputs(node.DynamicInputs(node.In, node.TypeFloat, nil, count, buildCount.Min, buildCount.Max, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, settings)
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
