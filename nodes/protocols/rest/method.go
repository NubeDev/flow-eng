package rest

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type HTTP struct {
	*node.Spec
	client    *resty.Client
	lastValue interface{}
}

func NewHttpWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, httpNode, category)
	url := node.BuildInput(node.URL, node.TypeString, nil, body.Inputs, false, false)
	reqBody := node.BuildInput(node.Body, node.TypeString, nil, body.Inputs, false, false)
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs, false, false)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, false, false)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, false)

	inputs := node.BuildInputs(url, reqBody, filter, trigger, enable)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)

	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	n := &HTTP{body, resty.New(), nil}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *HTTP) request(method string, body string) (resp *resty.Response, err error) {
	resp = &resty.Response{}
	client := inst.getClient()
	var parsedBody interface{}
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(body), &jsonMap)
	if err != nil {
		fmt.Println(err, "parse fail of type map")
	} else {
		parsedBody = jsonMap
	}

	err = json.Unmarshal([]byte(body), &parsedBody)
	if err != nil {
		fmt.Println(err, "parse fail of type interface{}")
	} else {
		parsedBody = jsonMap
	}

	url, _ := inst.ReadPinAsString(node.URL)
	if method == patch {
		resp, err = client.R().
			SetBody(parsedBody).
			Patch(url)
		return resp, err
	}
	if method == get {
		resp, err = client.R().
			Get(url)
		return resp, err
	}

	return resp, err
}

func (inst *HTTP) do() {
	filter, _ := inst.ReadPinAsString(node.Filter)
	reqBody, _ := inst.ReadPinAsString(node.Body)
	method, err := inst.getSettings()
	if err != nil {
		log.Error(err)
		return
	}
	resp, err := inst.request(method.Method, reqBody)
	if err != nil {
		inst.WritePin(node.Out, err.Error())
		return
	}
	if filter != "" {
		value := gjson.Get(resp.String(), filter)
		inst.lastValue = value.String()
		inst.WritePin(node.Out, value.String())
	} else {
		inst.lastValue = resp.String()
		inst.WritePin(node.Out, resp.String())
	}

}

func (inst *HTTP) getClient() *resty.Client {
	return inst.client
}

func (inst *HTTP) Process() {
	_, cov := inst.InputUpdated(node.TriggerInput)
	if cov {
		go inst.do()
	} else {
		inst.WritePin(node.Out, inst.lastValue)
	}

}
