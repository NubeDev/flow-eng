package math

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
)

type Add struct {
	*node.Spec
}

func NewAdd(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, add, category)
	if err != nil {
		return nil, err
	}
	return &Add{body}, nil
}

func getPoints() {
	client := resty.New()
	resp, err := client.R().
		SetResult(&node.Spec{}).
		Get("http://192.168.15.190:1660/api/points")
	fmt.Println(err)
	fmt.Println(resp.Status())
	// fmt.Println(resp.String())
}

func (inst *Add) Process() {
	process(inst)
	//go getPoints()
}

func (inst *Add) Cleanup() {}
