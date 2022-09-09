package timing

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/jobs"
	"github.com/NubeDev/flow-eng/node"
	"github.com/go-co-op/gocron"
	"time"
)

type Inject struct {
	*node.BaseNode
	cron         *gocron.Scheduler
	triggered    bool
	acknowledged bool

	lastTrigger string
	lastTime    time.Time
}

func NewInject(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body)
	body.Info.Name = inject
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, body.Outputs))
	j := jobs.New().Get()
	j.StartAsync()
	return &Inject{body, j, false, false, "", time.Now()}, nil
}

var count int

func job() {
	fmt.Println("*****************run job")
	count++
}

func (inst *Inject) check() {
	if !inst.triggered {
		inst.cron.Every(5).Second().Do(job)
	}
	inst.triggered = true

}

func (inst *Inject) Process() {

	inst.check()

	_, in1Val, _ := inst.ReadPinNum(node.In1)
	inst.WritePinNum(node.Out1, in1Val)

	fmt.Println("&&&&&&&&&&job count", count)
	fmt.Println("job trigger odd", count%2 == 0)
	fmt.Println("job trigger", "even", count%2 != 0)
}

func (inst *Inject) Cleanup() {}
