package inject

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/trigger"
)

type Inject struct {
	*node.Spec
	s map[string]interface{}
}

func NewInject(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, trigger.Inject, trigger.Category)
	// TODO: input only accepts strings currently, where by definition, input supports all data type
	message := node.BuildInput(node.Message, node.TypeString, nil, body.Inputs, false, false)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, false)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(message, enable, trigger)

	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	s := body.GetSettings()
	body.SetHelp("When ‘enable’ is ‘true’ and `trigger` transitions from `false` to `true`, `message` value is passed to `output`.  `message` and `output` values can be any Wires data types.")

	return &Inject{body, s}, nil
}

func (inst *Inject) Process() {
	message := inst.ReadPin(node.Message)
	enable, enableNull := inst.ReadPinAsBool(node.Enable)
	_, covBool, _ := inst.InputUpdated(node.TriggerInput)

	// fall back to values set in setting if input is not connected
	if enableNull && inst.s["enable"] != nil {
		if inst.s["enable"].(string) == "true" {
			enable = true
		} else {
			enable = false
		}
	}

	if enable && covBool {
		if message != nil {
			inst.WritePin(node.Out, message)
		} else {
			inst.WritePin(node.Out, inst.s["message"])
		}
	}
}
