package count

import (
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type Ramp struct {
	*node.Spec
	count     float64
	breakLoop bool
	lock      bool
}

func NewRamp(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, rampNode, category)
	duration := node.BuildInput(node.Duration, node.TypeFloat, nil, body.Inputs, nil)
	min := node.BuildInput(node.MinInput, node.TypeFloat, nil, body.Inputs, nil)
	max := node.BuildInput(node.MaxInput, node.TypeFloat, nil, body.Inputs, nil)
	body.Inputs = node.BuildInputs(duration, min, max)
	out := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(out)
	return &Ramp{body, 0, false, false}, nil
}

func (inst *Ramp) run() {
	duration, _ := inst.ReadPinAsDuration(node.Duration)
	min, _ := inst.ReadPinAsFloat(node.MinInput)
	max, _ := inst.ReadPinAsFloat(node.MaxInput)
	if !inst.breakLoop {
		inst.count += 1
	}
	if inst.count >= max {
		inst.breakLoop = true
	}
	if inst.breakLoop {
		inst.count -= 1
	}
	if inst.count == min && inst.breakLoop {
		inst.breakLoop = false
	}
	inst.lock = true
	time.Sleep(duration * time.Second)
	inst.lock = false
}

func (inst *Ramp) Process() {
	if !inst.lock {
		go inst.run()
		inst.WritePin(node.CountOut, inst.count)
	} else {
		inst.WritePin(node.CountOut, inst.count)
	}

}
