package flowctrl

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/graph"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/uuid"
	log "github.com/sirupsen/logrus"
)

type Flow struct {
	Graphs []*graph.Ordered
	nodes  []node.Node
}

func New(nodes ...node.Node) *Flow {
	runners := makeRunners(nodes)
	ordered := makeGraphs(runners)
	return &Flow{ordered, nodes}
}

func (p *Flow) Get() *Flow {
	return p
}

func (p *Flow) GetNodes() []node.Node {
	return p.nodes
}

func (p *Flow) GetNodesBase() []*node.BaseNode {
	var list []*node.BaseNode
	for _, n := range p.GetNodes() {
		n_ := &node.BaseNode{
			Inputs:   n.GetInputs(),
			Outputs:  n.GetOutputs(),
			Info:     n.GetInfo(),
			Settings: n.GetSettings(),
			Metadata: n.GetMetadata(),
		}
		list = append(list, n_)
	}
	return list
}

func (p *Flow) GetNodeBase(id string) *node.BaseNode {
	n := p.GetNode(id)
	return &node.BaseNode{
		Inputs:   n.GetInputs(),
		Outputs:  n.GetOutputs(),
		Info:     n.GetInfo(),
		Settings: n.GetSettings(),
		Metadata: n.GetMetadata(),
	}
}

func (p *Flow) AddNode(node node.Node) *Flow {
	flows := p.Get()
	flows.nodes = append(flows.nodes, node)
	runners := makeRunners(flows.nodes)
	ordered := makeGraphs(runners)
	flows.Graphs = ordered
	return flows
}

// ReBuildFlow makes all the node connections
func (p *Flow) ReBuildFlow(makeConnection bool) {
	for _, n := range p.GetNodes() {
		err := p.nodeConnector(n.GetID(), makeConnection)
		if err != nil {
			log.Errorf(fmt.Sprintf("rebuild-flow: node-id:%s node-name:%s err%s", n.GetID(), n.GetID(), err.Error()))
		}
	}
	for _, n := range p.GetNodes() {
		p.rebuildNode(n)
	}
}

func (p *Flow) GetNode(id string) node.Node {
	for _, n := range p.Get().nodes {
		if n.GetID() == id {
			return n
		}
	}
	return nil
}

func (p *Flow) GetNodeRunner(id string) *node.Runner {
	for _, n := range p.Get().Graphs {
		for _, runner := range n.Runners {
			if runner.NodeId() == id {
				return runner
			}
		}
	}
	return nil
}

func RemoveIndex(s []node.Node, index int) []node.Node {
	ret := make([]node.Node, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func (p *Flow) rebuildNode(node node.Node) {
	for i, n := range p.Get().nodes {
		if n.GetID() == node.GetID() {
			p.nodes = RemoveIndex(p.nodes, i)
		}
	}
	p.AddNode(node)
}

// ManualNodeConnector this node is just to really be used at the moment when testing the framework by making building the flow through code and not the json import
func (p *Flow) ManualNodeConnector(sourceNode, destNode node.Node, outPort, inPort node.PortName) error {
	for _, output := range sourceNode.GetOutputs() {
		if output.Name == outPort {
			for _, input := range destNode.GetInputs() {
				if input.Name == inPort {
					output.Connect(input.InputPort)
					log.Infof("manual-connection: source-node:%s dest-node:%s source-port:%s dest-port:%s", sourceNode.GetNodeName(), destNode.GetNodeName(), outPort, inPort)
					input.Connection.NodeID = sourceNode.GetID()
					input.Connection.NodePort = outPort
					return nil // connection was made so return
				}
			}
		}
	}
	return errors.New(fmt.Sprintf("failed to connect source-node:%s dest-node:%s source-port:%s dest-port:%s", sourceNode.GetNodeName(), destNode.GetNodeName(), outPort, inPort))

}

// nodeConnector will make the connections from nodeA to nodeA
//	-makeConnection if false will not make the connection, this would be set to false only when you want a snapshot of the current flow
func (p *Flow) nodeConnector(sourceID string, makeConnection bool) error {
	sourceNode := p.GetNode(sourceID)
	if sourceNode == nil {
		return errors.New("failed to find node by that id")
	}
	for _, output := range sourceNode.GetOutputs() {
		for _, connection := range output.Connections {
			destID := connection.NodeID
			if destID == "" {
				continue
			}
			destNode := p.GetNode(destID)
			if destNode == nil {
				return errors.New("failed to match ports for node connections")
			}
			for _, input := range destNode.GetInputs() {
				inPort := input.Name
				if inPort == connection.NodePort {
					if sourceID == input.Connection.NodeID {
						log.Infof("source-node:%s dest-node:%s source-port:%s dest-port:%s", sourceNode.GetNodeName(), destNode.GetNodeName(), output.Name, inPort)
						if makeConnection {
							output.OutputPort.Connect(input.InputPort)
						}
						input.Connection.NodePort = output.Name
						input.Connection.NodeID = sourceNode.GetID()
						for _, connection := range output.Connections {
							connection.NodeID = destNode.GetID()
							connection.NodePort = inPort
						}
					}
				}
			}
		}
	}
	return nil
}

func makeRunners(nodes []node.Node) []*node.Runner {
	nodesCount := len(nodes)
	runners := make([]*node.Runner, 0, nodesCount)
	for i := 0; i < nodesCount; i++ {
		n := nodes[i]

		runner := node.NewRunner(n)
		runners = append(runners, runner)
	}
	return runners
}

func makeGraphs(runners []*node.Runner) []*graph.Ordered {
	graphs := make([]*graph.Ordered, 0, 1)
	mapped := mapRunners(runners)
	// find root nodes and divide nodes by graphs
	for i := 0; i < len(runners); i++ {
		runner := runners[i]
		// root nodes are nodes with no connectors attached to output ports
		if isRootNode(runner) {
			graphs = append(graphs, graph.NewOrdered(runner, mapped))
		}
	}

	if len(graphs) == 0 {
		// panic("circular flows are not supported")
	}

	return graphs
}

func isRootNode(runner *node.Runner) bool {
	outputs := runner.Outputs()
	outputsConnected := false
	for i := 0; i < len(outputs); i++ {
		output := outputs[i]
		if len(output.Connectors()) > 0 {
			outputsConnected = true
		}
	}
	return !outputsConnected
}

func mapRunners(runners []*node.Runner) map[uuid.Value]*node.Runner {
	mapped := make(map[uuid.Value]*node.Runner)

	// map output port UUIDs to runners
	for i := 0; i < len(runners); i++ {
		runner := runners[i]
		outputs := runner.Outputs()
		for j := 0; j < len(outputs); j++ {
			output := outputs[j]
			mapped[output.UUID()] = runner
		}
	}
	return mapped
}
