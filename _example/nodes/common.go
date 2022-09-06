package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
)

type NodeSpec struct {
	NodeType    string
	InputCount  int
	OutputCount int
}

func Check(body *Node, nodeSpec NodeSpec) (*Node, error) {
	if body == nil {
		return nil, errors.New("node body can not be empty")
	}
	if body.Info.Name == "" {
		body.Info.Name = helpers.RandomName("node")
	}
	if body.Info.NodeID == "" {
		body.Info.NodeID = helpers.ShortUUID(nodeSpec.NodeType)
	}
	if len(body.InputList) != nodeSpec.InputCount {
		return nil, errors.New(fmt.Sprintf("input count is incorrect required:%d provided:%d", nodeSpec.InputCount, len(body.InputList)))
	}
	if len(body.OutputList) != nodeSpec.OutputCount {
		return nil, errors.New(fmt.Sprintf("output count is incorrect required:%d provided:%d", nodeSpec.OutputCount, len(body.OutputList)))
	}
	return body, nil
}
