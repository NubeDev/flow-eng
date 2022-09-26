package conversion

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
)

type NumToString struct {
	*node.Spec
}

func NewNumToString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, numToString, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &NumToString{body}, nil
}

func (inst *NumToString) Process() {
	in1 := inst.ReadPinAsFloat(node.In)
	v, ok := conversions.GetFloatOk(in1)
	if ok {
		inst.WritePin(node.Out, fmt.Sprintf("%f", v))
	} else {
		inst.WritePin(node.Out, nil)
	}
}

func (inst *NumToString) Cleanup() {}
