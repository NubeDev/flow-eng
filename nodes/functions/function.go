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
	eng        *rules.RuleEngine
	lastResult string
	lock       bool
	lockCount  int
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
	in1 := node.BuildInput(node.In1, node.TypeString, nil, body.Inputs, false, false)
	in2 := node.BuildInput(node.In2, node.TypeString, nil, body.Inputs, false, false)
	enable := node.BuildInput(node.Enable, node.TypeBool, true, body.Inputs, false, false)
	onlyRunOnStart := node.BuildInput(node.RunOnStartOnce, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(in1, in2, enable, onlyRunOnStart)
	output := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	outMsg := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	errOut := node.BuildOutput(node.ErrMsg, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(output, outMsg, errOut)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	eng := rules.NewRuleEngine()
	return &Func{body, eng, "", false, 0}, nil
}

/*
JSON parse
----------------------------------
let out = JSON.parse(Input.in1)
RQL.Result = out/10
----------------------------------


JSON stringify
----------------------------------
let out = {
	"_16":Input.in1*10
}
RQL.Result =  JSON.stringify(out)
----------------------------------

*/

func (inst *Func) Process() {
	if inst.lock {
		inst.lockCount++
		inst.writeValues(nil, rules.Processing)
		return
	} else {
		go inst.process()
		inst.lockNode(true)
	}

}

func (inst *Func) process() {

	if inst.disable() {
		inst.writeValues(nil, rules.Disabled)
		return
	}
	if inst.allowToRunFirstLoop() { // only execute on the first loop
	} else {
		inst.writeValues(nil, rules.Completed)
		return
	}

	code, err := getSettings(inst.Settings)
	if err != nil {
		inst.writeValues(err, rules.Error)
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
	updatedIn1, _, in1 := inst.InputUpdated(node.In1)
	updatedIn2, _, in2 := inst.InputUpdated(node.In2)

	if !updatedIn1 && !updatedIn2 { // write the last value
		inst.writeValues(err, rules.InputValuesNotUpdated)
		return
	}

	nodeInputs := map[string]interface{}{"in1": in1, "in2": in2}
	props["Input"] = nodeInputs

	err = inst.eng.AddRule(rule, props)
	if err != nil {
		inst.writeValues(err, rules.Error)
		return
	}

	res, err := inst.eng.ExecuteAndRemove(name, props, true)
	if err != nil {
		inst.writeValues(err, rules.Error)
		return
	}
	inst.lastResult = res.String()
	inst.writeValues(nil, rules.Completed)
	inst.lockCount = 0

}

func (inst *Func) writeValues(err error, state rules.State) {
	if err != nil {
		inst.WritePin(node.ErrMsg, err.Error())
	} else {
		inst.WritePin(node.ErrMsg, "")
	}
	inst.WritePin(node.Out, inst.lastResult)
	if inst.lockCount > 0 {
		inst.WritePin(node.Msg, fmt.Sprintf("%s %d", state, inst.lockCount))
	} else {
		inst.WritePin(node.Msg, state)
	}

	if state == rules.Completed {
		inst.lockNode(false)
	} else {
		inst.lockNode(true)
	}
}

func (inst *Func) lockNode(state bool) {
	inst.lock = state
}

func (inst *Func) disable() bool {
	enable, null := inst.ReadPinAsBool(node.Enable)
	if null { // can run
		return false
	}
	if enable { // can run
		return false
	} else { // disabled
		return true
	}

}

func (inst *Func) allowToRunFirstLoop() bool {
	_, firstLoop := inst.Loop()
	runOnStart, _ := inst.ReadPinAsBool(node.RunOnStartOnce) // only run on start

	if !runOnStart { // is disabled so pass
		return true
	}

	if runOnStart && firstLoop { // allow to run
		return true
	} else {
		return false // disable
	}

}
