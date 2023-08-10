package mathematics

import (
	"github.com/NubeDev/flow-eng/node"
)

type Advanced struct {
	*node.Spec
}

func NewAdvanced(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, mathAdvanced, Category)
	if err != nil {
		return nil, err
	}
	return &Advanced{body}, nil
}

func (inst *Advanced) Process() {
	process(inst)
}
