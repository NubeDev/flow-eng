package rest

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type HttpWrite struct {
	*node.Spec
	*resty.Client
}

func NewHttpWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, writeNode, category)
	url := node.BuildInput(node.URL, node.TypeString, nil, body.Inputs)
	reqBody := node.BuildInput(node.Body, node.TypeString, nil, body.Inputs)
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)

	inputs := node.BuildInputs(url, reqBody, filter, trigger, enable)
	out := node.BuildOutput(node.Result, node.TypeString, nil, body.Outputs)

	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &HttpWrite{body, resty.New()}, nil
}

func (inst *HttpWrite) request(method, body interface{}) (resp *resty.Response, err error) {
	resp = &resty.Response{}
	client := inst.getClient()
	url := inst.ReadPinAsString(node.URL)
	fmt.Println(url)
	if method == patch {
		resp, err = client.R().
			EnableTrace().
			SetBody(body).
			Patch(url)
		fmt.Println(3333, err)
		fmt.Println(resp.Status(), resp.String())
		return resp, err
	}

	return resp, err
}

func (inst *HttpWrite) do() {
	filter := inst.ReadPinAsString(node.Filter)
	reqBody := inst.ReadPinAsString(node.Body)
	fmt.Println(reqBody)
	method, err := getSettings(inst.GetSettings())
	if method == "" {
		method = patch
	}

	resp, err := inst.request(method, reqBody)
	fmt.Println(method, err)
	if err != nil {
		inst.WritePin(node.Result, nil)
		return
	}
	fmt.Println(resp.Status(), err)
	if filter != "" {
		value := gjson.Get(resp.String(), filter)
		inst.WritePin(node.Result, value.String())
	} else {
		inst.WritePin(node.Result, resp.String())
	}

}

func (inst *HttpWrite) getClient() *resty.Client {
	return inst.Client
}

func (inst *HttpWrite) Process() {
	_, cov := inst.InputUpdated(node.TriggerInput)
	if cov {
		go inst.do()

	} else {
		inst.WritePin(node.Result, nil)
	}

}

func (inst *HttpWrite) Cleanup() {}
