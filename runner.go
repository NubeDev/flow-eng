package flowctrl

import (
	"sync"
)

type SerialRunner struct {
	flow *Flow
}

func NewSerialRunner(flow *Flow) *SerialRunner {
	return &SerialRunner{flow}
}

func (sr *SerialRunner) Start() {
	for _, graph := range sr.flow.Graphs {
		for _, runner := range graph.Runners {
			runner.Start()
		}
	}
}

func (sr *SerialRunner) Process() {
	for _, graph := range sr.flow.Graphs {
		for _, runner := range graph.Runners {
			runner.Process()
		}
	}
	for _, graph := range sr.flow.Graphs {
		for _, runner := range graph.Runners {
			runner.Reset()
		}
	}
	return
}

func (sr *SerialRunner) Stop() {
	var wg sync.WaitGroup // stop function takes long time, and it gets accumulated without having wait group
	wg.Add(len(sr.flow.nodes))
	for _, graph := range sr.flow.Graphs {
		for _, runner := range graph.Runners {
			go runner.Stop(&wg)
		}
	}
	wg.Wait()
}
