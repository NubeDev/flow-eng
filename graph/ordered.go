package graph

import (
	"github.com/NubeDev/flow-eng/helpers/uuid"
	"github.com/NubeDev/flow-eng/node"
)

type Ordered struct {
	root    *node.Runner
	Runners []*node.Runner
}

func NewOrdered(root *node.Runner, mapped map[uuid.Value]*node.Runner) *Ordered {
	graph := &Ordered{root, make([]*node.Runner, 0, 1)}
	graph.build(mapped)
	return graph
}

func (ordered *Ordered) build(mapped map[uuid.Value]*node.Runner) {
	addedMap := make(map[uuid.Value]bool)
	ordered.buildRecursive(ordered.root, mapped, addedMap)
}

func (ordered *Ordered) buildRecursive(node *node.Runner, mapped map[uuid.Value]*node.Runner, added map[uuid.Value]bool) {
	connectors := node.Connectors()
	for i := 0; i < len(connectors); i++ {
		baseRunner, ok := mapped[connectors[i].FromUUID()]
		if !ok {
			panic("failed to find runner by output ID")
		}
		ordered.buildRecursive(baseRunner, mapped, added)
	}
	// add node only if it wasn't already added
	if _, ok := added[node.UUID()]; !ok {
		ordered.Runners = append(ordered.Runners, node)
		added[node.UUID()] = true
	}
}
