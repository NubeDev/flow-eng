package flowctrl

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/graph"
	"github.com/NubeDev/flow-eng/helpers/uuid"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

type Flow struct {
	Graphs []*graph.Ordered
	nodes  []node.Node
}

func New() *Flow {
	return &Flow{nil, nil}
}

func (p *Flow) Get() *Flow {
	return p
}

func (p *Flow) GetNodes() []node.Node {
	return p.nodes
}

type Message struct {
	Message string `json:"message"`
}

// SetNodePayload write value to a node from an api
func (p *Flow) SetNodePayload(uuid string, payload *node.Payload) (*Message, error) {
	for _, n := range p.GetNodes() {
		if uuid == n.GetID() {
			n.SetPayload(payload)
			return &Message{Message: "ok"}, nil
		}
	}
	return nil, errors.New("node not found")
}

// NodesValue get a single node value
func (p *Flow) NodesValue(uuid string) (*node.Values, error) {
	for _, n := range p.GetNodes() {
		if uuid == n.GetID() {
			return n.NodeValues(), nil
		}
	}
	return nil, errors.New("node not found")
}

// NodesValues get all the node current values from the runtime
func (p *Flow) NodesValues() []*node.Values {
	var out []*node.Values
	for _, n := range p.GetNodes() {
		out = append(out, n.NodeValues())
	}
	return out
}

func (p *Flow) GetNodesSpec() []*node.Spec {
	var list []*node.Spec
	for _, n := range p.GetNodes() {
		n_ := node.ConvertToSpec(n)
		list = append(list, n_)
	}
	return list
}

func (p *Flow) GetNodeSpec(id string) *node.Spec {
	n := p.GetNode(id)
	return node.ConvertToSpec(n)
}

func (p *Flow) AddNodes(node ...node.Node) {
	for _, n := range node {
		log.Infof("add node: %s %s", n.GetID(), n.GetName())
		p.addNode(n)
	}
}

func (p *Flow) AddNode(node node.Node) *Flow {
	return p.addNode(node)
}

func (p *Flow) addNode(node node.Node) *Flow {
	flows := p.Get()
	flows.nodes = append(flows.nodes, node)
	return flows
}

func (p *Flow) MakeGraph() *Flow {
	flows := p.Get()
	runners := makeRunners(flows.nodes)
	ordered := makeGraphs(runners)
	flows.Graphs = ordered
	return flows
}

// MakeNodeConnections makes all the node connections
func (p *Flow) MakeNodeConnections(makeConnection bool) {
	for _, n := range p.GetNodes() {
		err := p.nodeConnector(n.GetID(), makeConnection)
		if err != nil {
			log.Errorf(fmt.Sprintf("rebuild-flow: node-id:%s node-name:%s err%s", n.GetID(), n.GetID(), err.Error()))
		}
	}
}

func removeNode(s []node.Node, index int) []node.Node {
	ret := make([]node.Node, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func (p *Flow) rebuildNode(node node.Node) {
	for i, n := range p.Get().nodes {
		if n.GetID() == node.GetID() {
			p.nodes = removeNode(p.nodes, i)
		}
	}
	p.addNode(node)
}

func (p *Flow) WipeFlow() *Flow {
	p.nodes = nil
	p.Graphs = nil
	return p
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

// nodeConnector will make the connections from nodeA to nodeA
// for example we will make a link from the math-const node to the math add-node
//	-makeConnection if false will not make the link, this would be set to false only when you want a snapshot of the current flow
func (p *Flow) nodeConnector(nodeId string, makeConnection bool) error {
	getNode := p.GetNode(nodeId) // this is the add node
	if getNode == nil {
		return errors.New(fmt.Sprintf("node-connector: failed to find node id:%s", nodeId))
	}
	if len(getNode.GetInputs()) > 0 {
		// for a node we need its input and see if it has a link, if so we need the uuid of the node its link to
		for _, input := range getNode.GetInputs() { // this is the inputs from the add node
			// check the input has link
			connectionOutputName := input.Connection.NodePort // on const node will be named:out
			connectionOutputId := input.Connection.NodeID     // const node nodeId
			if connectionOutputName != "" {
				outputNode := p.GetNode(connectionOutputId) //this is the const node
				if outputNode != nil {
					for _, output := range outputNode.GetOutputs() {
						if output.Name == connectionOutputName {
							if makeConnection {
								output.Connect(input)
							}
							log.Infof("make node connections: %s:%s -> %s:%s", outputNode.GetName(), output.Name, getNode.GetName(), input.Name)
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
