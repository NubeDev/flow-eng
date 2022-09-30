package rest

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
)

type Get struct {
	*node.Spec
}

func NewGet(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, getNode, category)
	url := node.BuildInput(node.URL, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeFloat, nil, body.Inputs)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)

	inputs := node.BuildInputs(url, trigger, enable)
	out := node.BuildOutput(node.Result, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Get{body}, nil
}

func (inst *Get) Process() {
	url := inst.ReadPinAsString(node.URL)
	enable := inst.ReadPinAsString(node.Enable)
	fmt.Println(enable)
	fmt.Println(url)
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(url)

	fmt.Println(err)
	fmt.Println(resp.String())
	fmt.Println(resp.String())
}

func (inst *Get) Cleanup() {}
