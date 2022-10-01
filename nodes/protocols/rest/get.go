package rest

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type Get struct {
	*node.Spec
}

func build(body *node.Spec) *node.Spec {
	// ins
	url := node.BuildInput(node.URL, node.TypeString, nil, body.Inputs)
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeFloat, nil, body.Inputs)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(url, filter, trigger, enable)
	// outs
	out := node.BuildOutput(node.Result, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	return node.BuildNode(body, inputs, outputs, nil)

}

func NewGet(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, getNode, category)
	body = build(body)
	body.SetSchema(buildSchema())
	return &Get{body}, nil
}

func (inst *Get) Process() {
	url := inst.ReadPinAsString(node.URL)
	filter := inst.ReadPinAsString(node.Filter)
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
	if filter != "" {
		fmt.Println(2222, filter)
		//value := gjson.ParseBytes(resp.Body()).Get(filter)
		value := gjson.Get(resp.String(), filter)
		fmt.Println(value.String(), 11111)
		inst.WritePin(node.Result, value.String())
	} else {
		inst.WritePin(node.Result, resp.String())
	}

}

func (inst *Get) Cleanup() {}
