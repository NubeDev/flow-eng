package graph

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NubeDev/flow-eng/uuid"

	"github.com/NubeDev/flow-eng/buffer"

	"github.com/NubeDev/flow-eng/node"
)

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

type testNode struct {
	nodeInfo node.NodeInfo

	InputA  *node.InputPort
	InputB  *node.InputPort
	OutputA *node.OutputPort
	OutputB *node.OutputPort
}

func newTestNode() *testNode {
	inputa := node.NewInputPort(buffer.Int8)
	inputb := node.NewInputPort(buffer.Int8)
	outputa := node.NewOutputPort(buffer.Int8)
	outputb := node.NewOutputPort(buffer.Int8)
	info := node.NodeInfo{Name: "testNode", Description: "_", Version: "1.0.0"}
	return &testNode{info, inputa, inputb, outputa, outputb}
}

func (node *testNode) Info() node.NodeInfo {
	return node.nodeInfo
}

func (node *testNode) Process() {
	node.OutputA.Write([]byte{0})
	node.OutputB.Write([]byte{1})
}

func (node *testNode) Cleanup() {}

func TestOrdering_RunnersShuldNotBeDuplicated(t *testing.T) {
	// GIVEN
	nodeA := newTestNode()
	nodeB := newTestNode()
	nodeA.OutputA.Connect(nodeB.InputA)
	nodeA.OutputB.Connect(nodeB.InputB)
	nodeARunner := node.NewRunner(nodeA)
	nodeBRunner := node.NewRunner(nodeB)
	runnersMap := mapRunners([]*node.Runner{nodeARunner, nodeBRunner})
	// WHEN
	ordering := NewOrdered(nodeBRunner, runnersMap)
	// THEN
	for i := range ordering.Runners {
		runner := ordering.Runners[i]
		log.Println(runner.Name())
	}
	assert.Equal(t, 2, len(ordering.Runners))
}
