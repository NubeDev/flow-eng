package rest

import (
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
	url := node.BuildInput(node.URL, node.TypeString, nil, body.Inputs, nil)
	reqBody := node.BuildInput(node.Body, node.TypeString, nil, body.Inputs, nil)
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs, nil)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, nil)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, nil)

	inputs := node.BuildInputs(url, reqBody, filter, trigger, enable)
	out := node.BuildOutput(node.Outp, node.TypeString, nil, body.Outputs)

	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &HttpWrite{body, resty.New()}, nil
}

func (inst *HttpWrite) request(method, body interface{}) (resp *resty.Response, err error) {
	resp = &resty.Response{}
	client := inst.getClient()
	url, _ := inst.ReadPinAsString(node.URL)
	if method == patch {
		resp, err = client.R().
			EnableTrace().
			SetBody(body).
			Patch(url)
		return resp, err
	}
	return resp, err
}

func (inst *HttpWrite) do() {
	filter, _ := inst.ReadPinAsString(node.Filter)
	reqBody, _ := inst.ReadPinAsString(node.Body)
	method, err := getSettings(inst.GetSettings())
	if method == "" {
		method = patch
	}
	resp, err := inst.request(method, reqBody)
	if err != nil {
		inst.WritePinNull(node.Outp)
		return
	}
	if filter != "" {
		value := gjson.Get(resp.String(), filter)
		inst.WritePin(node.Outp, value.String())
	} else {
		inst.WritePin(node.Outp, resp.String())
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
		inst.WritePinNull(node.Outp)
	}

}
