package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
)

const httpHelp = `
A node used for sending HTTP/HTTPS requests
Methods supported: GET, POST, PUT, PATCH, DELETE

In the body you can send

- url
- trigger
- body
- method
- headers: {
    "Content-Type": "application/json",
    "Token": "abc1234"
}
- filter: filter the output response  (see for more info: https://github.com/tidwall/gjson#path-syntax)


example body using the function node


let in1 = Number(input.in1)
let msg = {}

msg.body = {
    "priority": {
        "_16": in1
    }
}

msg.method = "get"
msg.url = "http://0.0.0.0:1660/api/points/write/pnt_62f12094bf1b4fb1"
msg.trigger = false
msg.filter = "name"

RQL.Result =  JSON.stringify(msg)

`

const responseHelp = `
Will output the HTTP response
- body as a string
- status code
- errors
`

type HTTP struct {
	*node.Spec
	client    *resty.Client
	lastValue interface{}
}

func NewHttpWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, httpNode, category)
	input := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp(node.InHelp))
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp(node.FilterHelp))
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, false, node.SetInputHelp(node.EnableHelp))

	inputs := node.BuildInputs(input, filter, enable)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs, node.SetOutputHelp(node.OutHelp))
	response := node.BuildOutput(node.Response, node.TypeString, nil, body.Outputs, node.SetOutputHelp(responseHelp))

	outputs := node.BuildOutputs(out, response)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(httpHelp)
	n := &HTTP{body, resty.New(), nil}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *HTTP) Process() {
	bodyString, null := inst.ReadPinAsString(node.In)
	enable, _ := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.WritePin(node.Out, inst.lastValue)
	}
	if !null {
		go inst.processReq(bodyString)
	} else {
		inst.WritePin(node.Out, inst.lastValue)
	}
}

func (inst *HTTP) processReq(bodyString string) {
	filter, _ := inst.ReadPinAsString(node.Filter)
	method, err := inst.getSettings()
	if err != nil {
		log.Error(err)
		return
	}
	resp, responseOut, filterFromBody, err := inst.request(method.Method, bodyString)
	inst.WritePin(node.Response, jsonToString(responseOut))

	if err != nil {
		inst.WritePin(node.Out, err.Error())
		return
	}
	if filterFromBody != "" {
		filter = filterFromBody
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

func jsonToString(body interface{}) string {
	marshal, err := json.Marshal(body)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func (inst *HTTP) request(method string, bodyString string) (*resty.Response, *responseOutput, string, error) {
	body, err := inst.filterBody(method, bodyString)
	out := &responseOutput{}
	if body == nil {
		if err != nil {
			out.Error = err.Error()
		}
		return nil, out, "", err
	}
	resp, err := inst.httpSelect(body)
	if err != nil {
		return nil, out, "", err
	}
	if resp == nil {
		out.Error = "response was empty"
		return nil, out, "", err
	}
	out.Code = resp.StatusCode()
	out.Response = resp.String()
	return resp, out, body.Filter, err
}

func (inst *HTTP) getClient() *resty.Client {
	return inst.client
}

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

func (inst *HTTP) responseHandler(body *Body, err error, statusCode int, resp *resty.Response) string {
	if err != nil {
		return fmt.Sprintf("url: %s method: %s status-code: %d err: %v", body.URL, body.Method, statusCode, err)
	}
	if !resp.IsSuccess() {
		return fmt.Sprintf("url: %s method: %s status-code: %d", body.URL, body.Method, statusCode)
	}
	return resp.Status()
}

type responseOutput struct {
	Code     int         `json:"code"`
	Response interface{} `json:"response"`
	Error    string      `json:"error"`
	Message  string      `json:"message"`
}
