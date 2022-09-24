package conversion

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"strconv"
)

type StringToNum struct {
	*node.Spec
}

func NewStringToNum(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, stringToNum, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &StringToNum{body}, nil
}

func (inst *StringToNum) Process() {
	in1 := inst.ReadPin(node.In)
	if s, err := strconv.ParseFloat(fmt.Sprintf("%v", in1), 64); err == nil {
		inst.WritePin(node.Out, s)
	} else {
		inst.WritePin(node.Out, nil)
	}

}

func (inst *StringToNum) Cleanup() {}
