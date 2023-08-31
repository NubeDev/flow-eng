package rql

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

const getHelp = `

# A node used for sending HTTP/HTTPS requests

`

const responseHelp = `
Will output the HTTP response
- body as a string
- status code
- errors
`

type Get struct {
	*node.Spec
	client    *resty.Client
	lastValue interface{}
	locked    bool
}

func NewRQLGet(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, rqlGet, Category)
	input := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp(node.InHelp))
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs, false, false, node.SetInputHelp(node.FilterHelp))
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, false, node.SetInputHelp(node.EnableHelp))

	inputs := node.BuildInputs(input, filter, enable)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs, node.SetOutputHelp(node.OutHelp))
	response := node.BuildOutput(node.Response, node.TypeString, nil, body.Outputs, node.SetOutputHelp(responseHelp))

	outputs := node.BuildOutputs(out, response)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp(getHelp)
	n := &Get{body, resty.New(), nil, false}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Get) Process() {
	enable, null := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.WritePin(node.Out, inst.lastValue)
	}
	if !null {
		if !inst.locked {
			go inst.processReq()
			inst.WritePin(node.Out, inst.lastValue)
		} else {
			inst.WritePin(node.Out, inst.lastValue)
		}

	} else {
		inst.WritePin(node.Out, inst.lastValue)
	}
}

func (inst *Get) processReq() {
	inst.locked = true
	body := &body{}
	settings, err := getSettings(inst.GetSettings())
	if err != nil {
		return
	}
	get, err := inst.httpGet(getRule(settings.Name), body)
	inst.locked = false

	response, err := getRuleResponse(get.Body())
	if err != nil {
		fmt.Println(err)
		return
	}

	results, err := filterResults(response, settings.Results)
	if err != nil {
		fmt.Println(err)
		return
	}
	res := jsonToString(results)
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

func getRule(nameUUID string) string {
	return fmt.Sprintf("%srules/%s", baseURL(), nameUUID)
}

func baseURL() string {
	return fmt.Sprintf("http://0.0.0.0:1660/api/modules/module-core-rql/")
}

func (inst *Get) httpCommon(body *body) *resty.Request {
	return inst.client.R().SetHeaders(body.Headers).SetBody(body.Body)
}

func (inst *Get) httpGet(url string, body *body) (*resty.Response, error) {
	return inst.httpCommon(body).
		Get(url)
}

type body struct {
	URL     string            `json:"url"`
	Trigger bool              `json:"trigger"`
	Body    interface{}       `json:"body"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Filter  string            `json:"filter"` // https://github.com/tidwall/gjson#path-syntax
}

func jsonToString(body interface{}) string {
	marshal, err := json.Marshal(body)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func getRuleResponse(data []byte) (*rule, error) {
	r := &rule{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func filterResults(data *rule, resultsFilter string) (any, error) {
	if data == nil {
		return nil, errors.New("")
	}
	results := data.Result
	if len(results) > 0 {

		if resultsFilter == resultLatest {
			return results[0], nil
		}
		if resultsFilter == resultOldest {
			return results[len(results)-1], nil
		}
		if resultsFilter == resultAll {
			return results, nil
		}
	}
	return nil, errors.New("")
}

type rule struct {
	UUID              string        `json:"uuid"`
	Name              string        `json:"name"`
	LatestRunDate     string        `json:"latest_run_date"`
	Script            string        `json:"script"`
	Schedule          string        `json:"schedule"`
	Enable            bool          `json:"enable"`
	ResultStorageSize int           `json:"result_storage_size"`
	Result            []rulesResult `json:"rulesResult"`
}

type rulesResult struct {
	Result    interface{} `json:"rulesResult"`
	Timestamp string      `json:"timestamp"`
	Time      time.Time   `json:"time"`
}
