package timing

import (
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"log"
	"time"
)

type Delay struct {
	*node.Spec
	timer flowctrl.TimedDelay
}

func NewDelay(body *node.Spec, timer flowctrl.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delay, category)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	return &Delay{body, timer}, nil
}

func (inst *Delay) Process() {
	log.Println("Delayed START")
	if !inst.timer.WaitFor(5 * time.Second) {
		return
	}
	log.Println("Delayed triggered")
	in1 := inst.ReadPin(node.In)
	inst.WritePin(node.Out, in1)
}

func (inst *Delay) Cleanup() {}
