package functions

import (
	"bytes"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeIO/lib-goja/js"
	"github.com/mitchellh/mapstructure"
)

type Func struct {
	*node.Spec
	code string
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
	return &Func{body, ""}, nil
}

func (inst *Func) Process() {
	code, err := getSettings(inst.Settings)
	if err != nil {
		return
	}
	in1 := inst.ReadPin(node.In1)
	in2 := inst.ReadPin(node.In2)
	f, err := runFunc(in1, in2, code)
	if err != nil {
		fmt.Println(err)
	}

	inst.WritePin(node.Out, f)
}

func runFunc(val1, val2 interface{}, code string) (interface{}, error) {
	script, err := js.New(js.NewScript(code))
	if err != nil {
		return 0, err
	}
	arg := map[string]interface{}{"in1": val1, "in2": val2, "in3": val2}
	consoleLogs := new(bytes.Buffer)
	f, err := js.NewEngine().Execute(script, arg, js.WithLogging(consoleLogs))
	fmt.Println(consoleLogs)

	return f, err

}

/*
// example
let pri = {
    "priority": {
			// parse the string to a num Number(arg["in1"])
            "_14": Number(arg["in1"]),
            "_15": Number(arg["in2"]),
            "_16": Number(arg["in3"])

        }
}
// need to stringify otherwise the node would output a map
return JSON.stringify(pri)
*/
