package nodes

import (
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes/math"
	"github.com/NubeDev/flow-eng/_example/nodes/timing"
	"github.com/NubeDev/flow-eng/node"
)

func Builder(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case add:
		return math.NewAdd(body)
	case sub:
		return math.NewSub(body)
	case delay:
		return timing.NewDelay(body, flowctrl.NewTimer())
	}

	return nil, errors.New(fmt.Sprintf("no nodes found with name:%s", body.GetName()))
}
