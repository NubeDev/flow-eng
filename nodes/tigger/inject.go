package trigger

import (
	"github.com/NubeDev/flow-eng/helpers/jobs"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-co-op/gocron"
)

const DefaultInjectSecond = 2

type Inject struct {
	*node.Spec
	cron             *gocron.Scheduler
	triggerCount     uint32
	previousInterval int
}

func NewInject(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inject, category)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, node.SetInputHelp(node.IntervalInputHelp))
	body.Inputs = node.BuildInputs(interval)
	trigger := node.BuildOutput(node.Toggle, node.TypeBool, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(trigger)
	j := jobs.New().Get()
	j.StartAsync()
	return &Inject{body, j, 0, 0}, nil
}

func (inst *Inject) injectJob() {
	inst.triggerCount++
}

func (inst *Inject) run(interval int) {
	if interval == 0 {
		interval = DefaultInjectSecond
	}
	inst.cron.Clear() // if it's runtime change, we need to clear the running jobs and start fresh
	inst.cron.Every(interval).Second().Do(inst.injectJob)
}

func (inst *Inject) Process() {
	interval, _ := inst.ReadPinAsInt(node.Interval)
	if inst.previousInterval != interval { // if input interval gets changed on runtime
		inst.previousInterval = interval
		inst.run(interval)
	}
	toggle := inst.triggerCount%2 == 0
	inst.WritePin(node.Toggle, toggle)
}

func (inst *Inject) Stop() {
	inst.cron.Clear()
}
