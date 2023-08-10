package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Multiply struct {
	*node.Spec
}

func NewMultiply(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, multiply, Category)
	if err != nil {
		return nil, err
	}
	return &Multiply{body}, nil
}

func (inst *Multiply) Process() {
	process(inst)
}
