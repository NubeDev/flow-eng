package flowctrl

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/graph"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/uuid"
)

type Flow struct {
	Graphs []*graph.Ordered
	nodes  []*node.Node
}

func New(nodes ...*node.Node) *Flow {
	runners := makeRunners(nodes)
	ordered := makeGraphs(runners)
	return &Flow{ordered, nodes}
}

func (p *Flow) Get() *Flow {
	return p
}

func (p *Flow) GetNodes() []*node.Node {
	return p.nodes
}

func (p *Flow) GetNode(id string) *node.Node {
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

func (p *Flow) ReplaceNode(id string, node *node.Node) *node.Node {
	var found bool
	for _, n := range p.Get().nodes {
		if n.GetID() == id {
			n = nil
			found = true
		}
	}
	if found {
		p.AddNode(node)
	}
	return p.GetNode(id)
}

func (p *Flow) NodeConnector(sourceID string) error {
	a := p.GetNode(sourceID)
	if a == nil {
		return errors.New("failed to find node by that id")
	}
	for _, output := range a.Outputs {
		for _, connector := range output.Connections {
			destID := connector.NodeID
			if destID == "" {
				continue
			}
			fmt.Println("----------------")
			fmt.Printf(fmt.Sprintf("source-id:%s  dest-id:%s", sourceID, destID))
			fmt.Println()

			n := p.GetNode(destID)
			if n == nil {
				return errors.New("failed to match ports for node connections")
			}
			for _, input := range n.GetInputs() {
				port := input.Name

				if port == connector.NodePort {
					fmt.Println("*************************", port, connector.NodePort, sourceID, input.Connection.NodeID)
					fmt.Printf(fmt.Sprintf("source-id:%s  dest-port-id:%s", sourceID, input.Connection.NodeID))
					fmt.Println()
					if sourceID == input.Connection.NodeID {
						fmt.Println("*************************")
						output.OutputPort.Connect(input.InputPort)
					}
				}
			}
		}
	}
	return nil
}

func (p *Flow) AddNode(node *node.Node) *Flow {
	flows := p.Get()
	flows.nodes = append(flows.nodes, node)
	runners := makeRunners(flows.nodes)
	ordered := makeGraphs(runners)
	flows.Graphs = ordered
	return flows
}

func makeRunners(nodes []*node.Node) []*node.Runner {
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
		//panic("circular flows are not supported")
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
