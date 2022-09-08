package timing

import (
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"log"
	"time"
)

type Delay struct {
	*node.BaseNode
	Timer flowctrl.TimedDelay `json:"-"`
}

func NewDelay(body *node.BaseNode) (node.Node, error) {
	body = node.EmptyNode(body)
	body.Info.Name = delay
	body.Info.Category = category
	body.Info.NodeID = node.SetUUID(body.Info.NodeID)
	body.Inputs = node.BuildInputs(node.BuildInput(node.In1, node.TypeFloat, body.Inputs), node.BuildInput(node.In2, node.TypeFloat, body.Inputs))
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Out1, node.TypeFloat, body.Outputs))
	return &Delay{body, nil}, nil
}

func (inst *Delay) Process() {
	if !inst.Timer.WaitFor(5 * time.Second) {
		return
	}
	log.Println("Delayed triggered")
	_, in1Val, _ := inst.ReadPinNum(node.In1)
	inst.WritePinNum(node.Out1, in1Val)

}

func (inst *Delay) Cleanup() {}
