package subflow

import (
	"github.com/NubeDev/flow-eng/node"
)

type Loop struct {
	*node.Spec
}

func NewSubFlowFolder(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, folder, category)
	body.IsParent = true
	body = node.BuildNode(body, nil, nil, nil)
	return &Loop{body}, nil
}

func (inst *Loop) Process() {

}
