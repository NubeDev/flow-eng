package rql

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

const runHelp = `

# A node used for sending HTTP/HTTPS requests

`

type Run struct {
	*node.Spec
	client    *resty.Client
	lastValue interface{}
	locked    bool
}

func NewRQLRun(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, rqlRun, category)
	input := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp(node.InHelp))
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp(node.FilterHelp))
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, false, node.SetInputHelp(node.EnableHelp))

	inputs := node.BuildInputs(input, filter, enable)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs, node.SetOutputHelp(node.OutHelp))
	response := node.BuildOutput(node.Response, node.TypeString, nil, body.Outputs, node.SetOutputHelp(responseHelp))

	outputs := node.BuildOutputs(out, response)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(runHelp)
	n := &Run{body, resty.New(), nil, false}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Run) Process() {
	parsedBody, _ := inst.ReadPinAsString(node.In)
	enable, enableNull := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.WritePin(node.Out, inst.lastValue)
	}
	if !enableNull {
		if !inst.locked {
			go inst.processReq(parsedBody)
			inst.WritePin(node.Out, inst.lastValue)
		} else {
			inst.WritePin(node.Out, inst.lastValue)
		}

	} else {
		inst.WritePin(node.Out, inst.lastValue)
	}
}

func (inst *Run) processReq(bodyAsString string) {
	inst.locked = true
	inputBody, err := parseRuleBodyString(bodyAsString)
	if err != nil {
		return
	}
	settings, err := getSettings(inst.GetSettings())
	if err != nil {
		return
	}
	url := runRule(settings.Name)
	get, err := inst.httpPost(url, inputBody)
	inst.locked = false
	if get == nil {
		return
	}

	result, err := getRuleExistingResponse(get.Body())
	if err != nil {
		return
	}

	res := jsonToString(result.Return) // convert to string
	filter, _ := inst.ReadPinAsString(node.Filter)
	if filter != "" {
		value := gjson.Get(res, filter)
		inst.lastValue = value.String()
		inst.WritePin(node.Out, value.String())
	} else {
		inst.lastValue = res
		inst.WritePin(node.Out, res)
	}

}

func (inst *Run) httpCommon(body *runExistingBody) *resty.Request {
	return inst.client.R().SetBody(body)
}

func (inst *Run) httpPost(url string, body *runExistingBody) (*resty.Response, error) {
	var inputBody interface{}
	if body != nil {
		inputBody = body
	}
	return inst.httpCommon(body).
		SetBody(inputBody).
		Post(url)
}

type runExistingBody struct {
	Body interface{} `json:"body"`
}

type runExistingResult struct {
	Return    interface{} `json:"return"`
	Err       string      `json:"err"`
	TimeTaken string      `json:"time_taken"`
}

func getRuleExistingResponse(data []byte) (*runExistingResult, error) {
	r := &runExistingResult{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func parseRuleBodyString(bodyString string) (*runExistingBody, error) {
	b := &runExistingBody{}
	err := json.Unmarshal([]byte(bodyString), &b)
	return b, err
}

func runRule(nameUUID string) string {
	return fmt.Sprintf("%srules/run/%s", baseURL(), nameUUID)
}
