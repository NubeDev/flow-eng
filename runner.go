package flowctrl

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

type SerialRunner struct {
	flow *Flow
}

func NewSerialRunner(flow *Flow) *SerialRunner {
	return &SerialRunner{flow}
}

func (runner *SerialRunner) Process() (e error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			e = fmt.Errorf("flow processing error: %v", recovered)
		}
	}()

	for i := 0; i < len(runner.flow.Graphs); i++ {
		graph := runner.flow.Graphs[i]
		for j := 0; j < len(graph.Runners); j++ {
			runner := graph.Runners[j]
			err := runner.Process()
			if err != nil {
				// node was no triggered, not all input ports were written by dependent nodes
				if err == node.ErrNoInputData {
					continue
				}
				e = err
				return
			}
		}
		for j := 0; j < len(graph.Runners); j++ {
			runner := graph.Runners[j]
			runner.Reset()
		}
	}
	return
}
