package rest

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type HttpWrite struct {
	*node.Spec
}

func NewHttpWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, writeNode, category)
	url := node.BuildInput(node.URL, node.TypeString, nil, body.Inputs)
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeFloat, nil, body.Inputs)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)

	inputs := node.BuildInputs(url, filter, trigger, enable)
	out := node.BuildOutput(node.Result, node.TypeString, nil, body.Outputs)

	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, nil)
	return &HttpWrite{body}, nil
}

func (inst *HttpWrite) Process() {
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

func (inst *HttpWrite) Cleanup() {}
