package nodejson

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/tidwall/gjson"
)

type Filter struct {
	*node.Spec
}

func NewFilter(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, jsonFilter, category)
	in := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false)
	equation := node.BuildInput(node.Filter, node.TypeString, nil, body.Inputs, false)
	inputs := node.BuildInputs(in, equation)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Filter{body}, nil
}

func (inst *Filter) Process() {
	in1, _ := inst.ReadPinAsString(node.In)
	equation, _ := inst.ReadPinAsString(node.Filter)
	value := gjson.Get(in1, equation)
	inst.WritePin(node.Out, value.Value())
}
