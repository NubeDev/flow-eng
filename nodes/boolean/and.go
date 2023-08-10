package boolean

import (
	"github.com/NubeDev/flow-eng/node"
)

type And struct {
	*node.Spec
}

func NewAnd(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, and, Category)
	if err != nil {
		return nil, err
	}
	return &And{body}, nil
}

func (inst *And) Process() {
	Process(inst)
}
