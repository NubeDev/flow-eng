package functions

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/functions/rules"
	"github.com/NubeIO/module-core-rql/helpers/uuid"
	"github.com/mitchellh/mapstructure"
)

type Func struct {
	*node.Spec
	eng *rules.RuleEngine
}

type nodeSettings struct {
	Code string `json:"code"`
}

func getSettings(body map[string]interface{}) (string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return "", err
	}
	if settings != nil {
		return settings.Code, nil
	}
	return "", nil
}

func NewFunc(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, funcNode, category)
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeString, nil, 2, 3, 3, body.Inputs, false)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	eng := rules.NewRuleEngine()
	return &Func{body, eng}, nil
}

/*
JSON parse
----------------------------------
let out = JSON.parse(RQL.in1)
RQL.Result = out._16/10
----------------------------------


JSON stringify
----------------------------------
let out = {
	"_16":RQL.in1*10
}
RQL.Result =  JSON.stringify(out)
----------------------------------

*/

func (inst *Func) Process() {

	code, err := getSettings(inst.Settings)
	if err != nil {
		return
	}

	props := make(rules.PropertiesMap)
	props["Core"] = inst.eng

	name := uuid.ShortUUID()
	rule := &rules.RQL{
		Name:   name,
		Script: code,
		Enable: true,
	}
	props["RQL"] = rule
	in1 := inst.ReadPin(node.In1)
	in2 := inst.ReadPin(node.In2)

	nodeInputs := map[string]interface{}{"in1": in1, "in2": in2}
	props["RQLInput"] = nodeInputs
	err = inst.eng.AddRule(rule, props)
	if err != nil {
		return
	}

	res, err := inst.eng.ExecuteAndRemove(name, props, true)

	if err != nil {
		fmt.Println("ExecuteByName", err)
		// return
	}

	inst.WritePin(node.Out, res.String())
}
