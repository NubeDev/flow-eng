package rest

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
)

type HTTP struct {
	*node.Spec
	client    *resty.Client
	lastValue interface{}
}

func NewHttpWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, httpNode, category)
	input := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false)
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs, false, false)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, false)

	inputs := node.BuildInputs(input, filter, enable)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)

	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	n := &HTTP{body, resty.New(), nil}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *HTTP) Process() {
	reqBody, null := inst.ReadPinAsString(node.In)
	enable, _ := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.WritePin(node.Out, inst.lastValue)
	}
	if !null {
		go inst.processReq(reqBody)
	} else {
		inst.WritePin(node.Out, inst.lastValue)
	}
}

func (inst *HTTP) processReq(reqBody string) {
	filter, _ := inst.ReadPinAsString(node.Filter)
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

func (inst *HTTP) request(method string, bodyString string) (*resty.Response, error) {
	body, err := inst.filterBody(method, bodyString)
	if body == nil {
		return nil, err
	}
	resp, err := inst.httpSelect(body)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (inst *HTTP) getClient() *resty.Client {
	return inst.client
}

type bodyType string

const errNoBodyType = "http-node: no body type"
const errBodyEmpty = "http-node: body can not be empty"

func (inst *HTTP) filterBody(method, bodySting string) (*Body, error) {
	var err error
	body, err := parseBodyString(bodySting)
	if body == nil {
		return nil, errors.New(errBodyEmpty)
	}
	if err != nil {
		return nil, err
	}
	if method != "" {
		method = body.Method
	}
	body.Method = method
	return body, errors.New(errNoBodyType)
}

func parseBodyString(bodyString string) (*Body, error) {
	body := &Body{}
	err := json.Unmarshal([]byte(bodyString), &body)
	return body, err
}

type Body struct {
	URL     string            `json:"url"`
	Trigger bool              `json:"trigger"`
	Body    interface{}       `json:"body"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Filter  string            `json:"filter"` // https://github.com/tidwall/gjson#path-syntax
}

func (inst *HTTP) httpSelect(body *Body) (*resty.Response, error) {
	var method string
	var url string
	var err error
	if body.Method != "" {
		method = body.Method
	}
	url = body.URL
	method = strings.ToUpper(method)

	resp := &resty.Response{}
	switch method {
	case get:
		resp, err = inst.httpGet(url, body)
	case post:
		resp, err = inst.httpPost(url, body)
	case put:
		resp, err = inst.httpPut(url, body)
	case patch:
		resp, err = inst.httpPut(url, body)
	case httpDelete:
		resp, err = inst.httpDelete(url, body)
	}
	return resp, err
}

func (inst *HTTP) httpCommon(body *Body) *resty.Request {
	return inst.getClient().R().SetHeaders(body.Headers).SetBody(body.Body)
}

func (inst *HTTP) httpGet(url string, body *Body) (*resty.Response, error) {
	return inst.httpCommon(body).
		Get(url)
}

func (inst *HTTP) httpPost(url string, body *Body) (*resty.Response, error) {
	return inst.httpCommon(body).
		Post(url)
}

func (inst *HTTP) httpPut(url string, body *Body) (*resty.Response, error) {
	return inst.httpCommon(body).
		Put(url)
}

func (inst *HTTP) httpPatch(url string, body *Body) (*resty.Response, error) {
	return inst.httpCommon(body).
		Patch(url)
}

func (inst *HTTP) httpDelete(url string, body *Body) (*resty.Response, error) {
	return inst.httpCommon(body).
		Delete(url)
}
