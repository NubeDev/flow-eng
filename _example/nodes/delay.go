package nodes

import (
	"github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/node"
	"log"
	"time"
)

type Delay struct {
	info node.Info

	Input  *node.InputPort
	Output *node.OutputPort

	reader *adapter.Int8
	writer *adapter.Int8

	timer flowctrl.TimedDelay
}

func NewDelay(timer flowctrl.TimedDelay) *Delay {
	input := node.NewInputPort(buffer.Int8)
	output := node.NewOutputPort(buffer.Int8)
	info := node.Info{Name: "delay", Description: "passed data to another node with delay", Version: "1.0.0"}
	return &Delay{info, input, output, adapter.NewInt8(input), adapter.NewInt8(output), timer}
}

func (node *Delay) Info() node.Info {
	return node.info
}

func (node *Delay) Process() {
	if !node.timer.WaitFor(5 * time.Second) {
		return
	}
	log.Println("Delayed triggered")

	read := node.reader.Get()
	node.writer.Set(read)
}

func (node *Delay) Cleanup() {}
