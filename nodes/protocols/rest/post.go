package rest

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
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
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &HttpWrite{body}, nil
}

func (inst *HttpWrite) Process() {
	url := inst.ReadPinAsString(node.URL)
	filter := inst.ReadPin(node.Filter)
	fmt.Println(1111111)
	fmt.Println(filter)
	fmt.Println(1111111)
	//enable := inst.ReadPinAsString(node.Enable)

	client := resty.New()
	resp, err := client.R().
		SetBody(filter).
		Patch(url)

	fmt.Println(resp.Status())
	fmt.Println(err)
	fmt.Println(resp.String())
	if filter != "" {
		fmt.Println(2222, filter)
		//value := gjson.ParseBytes(resp.Body()).Get(filter)

		inst.WritePin(node.Result, 22)
	} else {
		inst.WritePin(node.Result, resp.String())
	}

}

func (inst *HttpWrite) Cleanup() {}
