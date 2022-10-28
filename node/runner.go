package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/global"
	"github.com/NubeDev/flow-eng/helpers/uuid"
	log "github.com/sirupsen/logrus"
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

func (runner *Runner) Process() error {
	// trigger all connectors to input ports
	err := runner.processConnectors()
	if err != nil {
		log.Errorf("RUNNER node:%s name-name%s err:%s", runner.node.GetNodeName(), runner.node.GetName(), err.Error())
		return err
	}
	// run processing node
	runner.node.Process()
	return nil
}

func (runner *Runner) Reset() {
	runner.resetConnectors()
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
			log.Errorf(fmt.Sprintf("err from runner:%s from:%s to:%s", err.Error(), conn.from.Name, conn.to.Name))
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
