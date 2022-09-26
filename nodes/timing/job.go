package timing

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/jobs"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-co-op/gocron"
	"time"
)

type Inject struct {
	*node.Spec
	cron        *gocron.Scheduler
	triggered   bool
	lastTrigger string
	lastTime    time.Time
}

func NewInject(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inject, category)

	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(interval)

	trigger := node.BuildOutput(node.Trigger, node.TypeFloat, nil, body.Outputs)
	toggle := node.BuildOutput(node.Toggle, node.TypeFloat, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(trigger, toggle)
	j := jobs.New().Get()
	j.StartAsync()
	return &Inject{body, j, false, "", time.Now()}, nil
}

var count int
var t bool

func set() {
	time.Sleep(2 * time.Second)
	t = false
}

func job() {
	count++
	t = true
	go set()

}

func (inst *Inject) check(interval interface{}) {
	if !inst.triggered {
		inst.cron.Every(interval).Second().Do(job)
	}
	inst.triggered = true

}

func (inst *Inject) Process() {
	interval := inst.ReadPinAsInt(node.Interval)
	inst.check(interval)
	inst.WritePin(node.Trigger, conversions.BoolToNum(t))           // set to on for 2 sec
	inst.WritePin(node.Toggle, conversions.BoolToNum(count%2 == 0)) // toggle on/off
}

func (inst *Inject) Cleanup() {}
