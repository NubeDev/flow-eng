package switches

import (
	"github.com/NubeDev/flow-eng/node"
)

type Switch struct {
	*node.Spec
}

func NewSwitch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, switchNode, category)
	inSwitch := node.BuildInput(node.Switch, node.TypeFloat, nil, body.Inputs) // TODO: this input shouldn't have a manual override value
	inTrue := node.BuildInput(node.InTrue, node.TypeFloat, nil, body.Inputs)
	inFalse := node.BuildInput(node.InFalse, node.TypeFloat, nil, body.Inputs)

	inputs := node.BuildInputs(inSwitch, inTrue, inFalse)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Switch{body}, nil
}

func (inst *Switch) Process() {
	inSwitch, _ := inst.ReadPinAsFloat(node.Switch)
	inTrue, _ := inst.ReadPinAsFloat(node.InTrue)
	inFalse, _ := inst.ReadPinAsFloat(node.InFalse)

	inSwitchAsBool := inSwitch == 1

	if inSwitchAsBool {
		inst.WritePin(node.Out, inTrue)
	} else {
		inst.WritePin(node.Out, inFalse)
	}
}
