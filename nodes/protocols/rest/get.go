package rest

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Get struct {
	*node.Spec
}

func build(body *node.Spec) *node.Spec {
	// ins
	url := node.BuildInput(node.URL, node.TypeString, nil, body.Inputs)
	filter := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs)
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
	return &Get{body}, nil
}

func (inst *Get) do() {
	url, _ := inst.ReadPinAsString(node.URL)
	filter, _ := inst.ReadPinAsString(node.Filter)
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(url)
	if err != nil {
		log.Error(err)
	}
	if filter != "" {
		value := gjson.Get(resp.String(), filter)
		inst.WritePin(node.Result, value.String())
	} else {
		inst.WritePin(node.Result, resp.String())
	}

}

func (inst *Get) Process() {
	_, cov := inst.InputUpdated(node.TriggerInput)
	if cov {
		go inst.do()

	} else {
		inst.WritePin(node.Result, nil)
	}

}
