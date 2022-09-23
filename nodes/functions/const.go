package functions

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeIO/lib-goja/js"
)

type Func struct {
	*node.Spec
}

func NewFunc(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, funcNode, category)
	inputs := node.BuildInputs(node.DynamicInputs(node.TypeFloat, nil, 2, 3, 3, body.Inputs)...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Func{body}, nil
}

func (inst *Func) Process() {
	in1 := inst.ReadPinAsFloat(node.In1)
	in2 := inst.ReadPinAsFloat(node.In2)
	f, err := runFunc(in1, in2)
	fmt.Println(err)
	inst.WritePin(node.Out, f)
}

func runFunc(val1, val2 float64) (float64, error) {
	code := `
	var out = in1*in2
	`
	j, err := js.New(code)

	if err != nil {
		fmt.Printf("Error loading JS code %v", err)
	}
	j.Set("in1", val1)
	j.Set("in2", val2)

	_, err = j.Run()
	if err != nil {
		fmt.Println(err)
	}
	res, err := j.GetGlobalObject().ToFloat("out")
	if err != nil {
		fmt.Println(err)
	}
	return res, nil

}

func (inst *Func) Cleanup() {}
