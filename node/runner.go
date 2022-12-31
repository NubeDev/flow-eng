package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/global"
	"github.com/NubeDev/flow-eng/helpers/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Runner struct {
	uuid       uuid.Value
	nodeId     string
	name       string
	node       Node
	inputs     []Port
	outputs    []Port
	connectors []*Connector
}

func NewRunner(node Node) *Runner {
	inputs := Ports(node, DirectionInput)
	outputs := Ports(node, DirectionOutput)
	connectors := Connectors(inputs)
	info := node.GetInfo()
	nodeID := node.GetID()
	id := uuid.New()
	name := fmt.Sprintf("%s_%d_%s", info.Name, id, nodeID)
	return &Runner{id, nodeID, name, node, inputs, outputs, connectors}
}

func (runner *Runner) Name() string {
	return runner.name
}

func (runner *Runner) NodeId() string {
	return runner.nodeId
}

func (runner *Runner) UUID() uuid.Value {
	return runner.uuid
}

func (runner *Runner) Start() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Errorf("flow on remove error: %v", recovered)
			log.Errorf("error on node: %s, node_id: %s", runner.node.GetName(), runner.NodeId())
		}
	}()

	// on multiple graphs the start lifecycle might be already got called for the same node
	if runner.node.GetCurrentState() != STARTED {
		runner.node.Start()
		runner.node.SetCurrentState(STARTED)
	}
}

func (runner *Runner) Process() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Errorf("flow process error: %v", recovered)
			log.Errorf("error on node: %s, node_id: %s", runner.node.GetName(), runner.NodeId())
		}
	}()

	// trigger all connectors to input ports
	err := runner.processConnectors()
	if err != nil {
		log.Errorf("RUNNER node: %s err: %s", runner.node.GetName(), err.Error())
		return
	}

	// if the same node is already processed, don't process it again
	// this restricts the same node from processing which has multiple output connections to the nodes
	if !runner.node.GetProcessed() {
		runner.node.Process()
		runner.node.SetProcessed()
	}
}

func (runner *Runner) Stop(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		if recovered := recover(); recovered != nil {
			log.Errorf("flow on start error: %v", recovered)
			log.Errorf("error on node: %s, node_id: %s", runner.node.GetName(), runner.NodeId())
		}
	}()

	// on multiple graphs the stop lifecycle might be already got called for the same node
	if runner.node.GetCurrentState() != STOPPED {
		runner.node.Stop()
		runner.node.SetCurrentState(STOPPED)
	}
}

func (runner *Runner) Reset() {
	runner.resetConnectors()
	runner.node.ResetProcessed()
}

func (runner *Runner) Outputs() []Port {
	return runner.outputs
}

func (runner *Runner) Inputs() []Port {
	return runner.inputs
}

func (runner *Runner) Connectors() []*Connector {
	return runner.connectors
}

func (runner *Runner) processConnectors() error {
	connectorsCount := len(runner.connectors)
	if connectorsCount == 0 {
		return nil
	}
	for i := 0; i < connectorsCount; i++ {
		conn := runner.connectors[i]
		debug := &global.Debug{}
		if global.DebugConnections {
			debug = &global.Debug{
				NodeUUID:   runner.NodeId(),
				NodeName:   runner.Name(),
				FromOutput: string(conn.from.Name),
				ToInput:    string(conn.to.Name),
			}
		}
		err := conn.Trigger(debug)
		if err != nil {
			log.Errorf("error from runner: %s from: %s to: %s", err.Error(), conn.from.Name, conn.to.Name)
			return err
		}
	}
	return nil
}

func (runner *Runner) resetConnectors() {
	connectorsCount := len(runner.connectors)
	if connectorsCount == 0 {
		return
	}
	for i := 0; i < connectorsCount; i++ {
		runner.connectors[i].Reset()
	}
}
