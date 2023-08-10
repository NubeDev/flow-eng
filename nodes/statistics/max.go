package statistics

import (
	"github.com/NubeDev/flow-eng/node"
)

type Max struct {
	*node.Spec
}

func NewMax(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, max, Category)
	if err != nil {
		return nil, err
	}
	return &Max{body}, nil
}

func (inst *Max) Process() {
	process(inst)
}
