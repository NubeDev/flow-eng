package statistics

import (
	"github.com/NubeDev/flow-eng/node"
)

type Min struct {
	*node.Spec
}

func NewMin(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, min, Category)
	if err != nil {
		return nil, err
	}
	return &Min{body}, nil
}

func (inst *Min) Process() {
	process(inst)
}
