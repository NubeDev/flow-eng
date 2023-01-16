package cov

import (
	"math"
	"time"

	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/trigger"
)

type COV struct {
	*node.Spec
	lastValue *float64
	running   bool
	s         map[string]interface{}
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
	s := body.GetSettings()
	body.SetHelp("when ‘input’ changes value, output becomes ‘true’ for ‘interval’ duration, then ‘output’ changes back to ‘false’. For Numeric ‘input’ values, the change of value must be greater than the ‘threshold’ value to trigger the output. Interval value must be equal or larger than 1.")

	return &COV{body, nil, false, s}, nil
}

func (inst *COV) Process() {
	var diff float64
	var covUnits interface{}

	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	covInterval, intervalNull := inst.ReadPinAsFloat(node.Interval)
	covThreshold, thresholdNull := inst.ReadPinAsFloat(node.Threshold)

	// fall back values in setting
	if thresholdNull && inst.s["covThreshold"] != nil {
		covThreshold = inst.s["covThreshold"].(float64)
	}
	if intervalNull && inst.s["interval"] != nil {
		covInterval = inst.s["interval"].(float64)
	}
	if inst.s["units"] == nil {
		covUnits = trigger.Seconds
	} else {
		covUnits = inst.s["units"]
	}

	// outputs false if the input is nil or there is no lastValue
	if inputNull || inst.lastValue == nil {
		inst.WritePinBool(node.Outp, false)
	} else {
		// call 'writeOutput' when the absolute diff between last two inputs are larger than 'covThreshold' and there are no previous routine running
		diff = math.Abs(input - *inst.lastValue)
		if diff > covThreshold && !inst.running {
			go writeOutput(inst, covInterval, covUnits)
			inst.running = true
		}
	}
	inst.lastValue = &input
}

func writeOutput(inst *COV, covInterval float64, covUnits interface{}) {
	var duration time.Duration
	switch covUnits.(string) {
	case string(trigger.Milliseconds):
		duration = time.Duration(covInterval * float64(time.Millisecond))
	case string(trigger.Seconds):
		duration = time.Duration(covInterval * float64(time.Second))
	case string(trigger.Minutes):
		duration = time.Duration(covInterval * float64(time.Minute))
	case string(trigger.Hours):
		duration = time.Duration(covInterval * float64(time.Hour))
	}

	// set the output pin to true for 'duration' period before setting it to false
	// set inst.running to false after routine is finished
	inst.WritePinBool(node.Outp, true)
	time.Sleep(duration)
	inst.WritePinBool(node.Outp, false)
	inst.running = false
}
