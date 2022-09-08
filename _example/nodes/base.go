package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/_example/nodes/math"
	"github.com/NubeDev/flow-eng/node"
)

func Builder(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case add:
		return math.NewAdd(body)
	case sub:
		return math.NewSub(body)
	}
	return nil, errors.New(fmt.Sprintf("no nodes found with name:%s", body.GetName()))
}
