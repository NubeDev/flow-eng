package switches

import (
	"github.com/NubeDev/flow-eng/node"
)

type Switch struct {
	*node.Spec
}

func NewSwitch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, switchNode, category)
	inSwitch := node.BuildInput(node.Switch, node.TypeBool, nil, body.Inputs, false, true) // TODO: this input shouldn't have a manual override value
	inTrue := node.BuildInput(node.InTrue, node.TypeFloat, nil, body.Inputs, false, false)
	inFalse := node.BuildInput(node.InFalse, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(inSwitch, inTrue, inFalse)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Switch{body}, nil
}

func (inst *Switch) Process() {
	inSwitch, _ := inst.ReadPinAsBool(node.Switch)
	inTrue, inTrueNull := inst.ReadPinAsFloat(node.InTrue)
	inFalse, inFalseNull := inst.ReadPinAsFloat(node.InFalse)

	if inSwitch {
		if inTrueNull {
			inst.WritePinNull(node.Out)
		} else {
			inst.WritePinFloat(node.Out, inTrue)
		}
	} else {
		if inFalseNull {
			inst.WritePinNull(node.Out)
		} else {
			inst.WritePinFloat(node.Out, inFalse)
		}
	}
}
