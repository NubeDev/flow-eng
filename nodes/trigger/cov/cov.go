package cov

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/trigger"
)

type COV struct {
	*node.Spec
	// Name         string
	// Interval     float64
	// Units        trigger.TimeUnits
	// COVThreshold int
}

func NewCOV(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, trigger.COV, trigger.Category)
	input := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs)
	threshold := node.BuildInput(node.Threshold, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(input, interval, threshold)
	out := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &COV{body}, nil
}

func (inst *COV) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		// settings, _ := getSettings(inst.GetSettings())
	}
	// min, _ := inst.ReadPinAsFloat(node.MinInput)
	// max, _ := inst.ReadPinAsFloat(node.MaxInput)
	// _, boolCov := inst.InputUpdated(node.TriggerInput)
	// if boolCov || firstLoop {
	// 	inst.value = float.RandFloat(min, max)
	// }
	// inst.WritePinFloat(node.Out, inst.value, inst.precision)
}
