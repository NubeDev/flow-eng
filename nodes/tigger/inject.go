package trigger

import (
	"github.com/NubeDev/flow-eng/helpers/jobs"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-co-op/gocron"
)

type Inject struct {
	*node.Spec
	cron         *gocron.Scheduler
	triggered    bool
	triggerCount uint64
}

func NewInject(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inject, category)

	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(interval)

	trigger := node.BuildOutput(node.Toggle, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(trigger)
	j := jobs.New().Get()
	j.StartAsync()
	return &Inject{body, j, false, 0}, nil
}

func (inst *Inject) injectJob() {
	inst.triggerCount++
}

func (inst *Inject) run(interval int) {
	if interval == 0 {
		interval = 2
	}
	inst.cron.Every(interval).Second().Do(inst.injectJob)
	inst.triggered = true
}

func (inst *Inject) Process() {
	interval := inst.ReadPinAsInt(node.Interval)
	if !inst.triggered {
		inst.run(interval)
	}
	toggle := inst.triggerCount%2 == 0
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.WritePin(node.Toggle, false)
	} else {
		inst.WritePin(node.Toggle, toggle)
	}

}

func (inst *Inject) Cleanup() {}
