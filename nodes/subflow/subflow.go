package subflow

import (
	"github.com/NubeDev/flow-eng/node"
)

type Folder struct {
	*node.Spec
}

func NewSubFlowFolder(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, folder, Category)
	body.IsParent = true
	body = node.BuildNode(body, nil, nil, nil)
	return &Folder{body}, nil
}

func (inst *Folder) Process() {

}
