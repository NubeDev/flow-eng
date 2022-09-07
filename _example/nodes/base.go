package nodes

import (
	"errors"
	"github.com/NubeDev/flow-eng/_example/nodes/math"
	"github.com/NubeDev/flow-eng/node"
)

func Builder(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case "add":
		return math.NewAdd(body)
	case "nodeB":
		//return NewNodeB(body)
	}
	return nil, errors.New("node not found")
}
