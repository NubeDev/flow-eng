package timing

import (
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"log"
	"time"
)

type Delay struct {
	*node.BaseNode
	timer flowctrl.TimedDelay
}

func NewDelay(body *node.BaseNode, timer flowctrl.TimedDelay) (node.Node, error) {
	body = node.Defaults(body, delay, category)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, nil, body.Outputs))
	return &Delay{body, timer}, nil
}

func (inst *Delay) Process() {
	log.Println("Delayed START")
	if !inst.timer.WaitFor(5 * time.Second) {
		return
	}
	log.Println("Delayed triggered")
	in1 := inst.ReadPin(node.In1)
	inst.WritePin(node.Out1, in1)
}

func (inst *Delay) Cleanup() {}
