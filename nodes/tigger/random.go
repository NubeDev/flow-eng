package trigger

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type Random struct {
	*node.Spec
	value     float64
	precision int
}

func NewRandom(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, randomFloat, category)
	min := node.BuildInput(node.Min, node.TypeFloat, nil, body.Inputs)
	max := node.BuildInput(node.Max, node.TypeFloat, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(min, max, trigger)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	precision, _ := getSettings(body.GetSettings())
	return &Random{body, 0, precision}, nil
}

func (inst *Random) Process() {
	_, firstLoop := inst.Loop()
	min, _ := inst.ReadPinAsFloat(node.Min)
	max, _ := inst.ReadPinAsFloat(node.Max)
	_, boolCov := inst.InputUpdated(node.TriggerInput)
	if boolCov || firstLoop {
		inst.value = float.RandFloat(min, max)
	}
	inst.WritePinFloat(node.Out, inst.value, inst.precision)
}
