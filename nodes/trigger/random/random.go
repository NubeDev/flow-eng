package random

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/trigger"
)

type Random struct {
	*node.Spec
	s map[string]interface{}
}

func NewRandom(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, trigger.RandomFloat, trigger.Category)
	min := node.BuildInput(node.MinInput, node.TypeFloat, nil, body.Inputs)
	max := node.BuildInput(node.MaxInput, node.TypeFloat, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(min, max, trigger)
	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	s := body.GetSettings()
	body.SetHelp("When ‘trigger’ transitions from ‘false’ to ‘true’, a random number between ‘min’ and ‘max’ values is produced at ‘output’. The number of decimal places that ‘output’ values have can be set from settings.")

	return &Random{body, s}, nil
}

func (inst *Random) Process() {
	min, minNull := inst.ReadPinAsFloat(node.MinInput)
	max, maxNull := inst.ReadPinAsFloat(node.MaxInput)
	_, covBool := inst.InputUpdated(node.TriggerInput)

	// fall back to values set in settings if no input
	if minNull && inst.s["min"] != nil {
		min = inst.s["min"].(float64)
	}

	if maxNull && inst.s["max"] != nil {
		max = inst.s["max"].(float64)
	}

	var p int
	if inst.s["percision"] != nil {
		p = int(inst.s["percision"].(float64))
	} else {
		p = 2
	}

	if covBool {
		inst.WritePinFloat(node.Outp, float.RandFloat(min, max), p)
	}

}
