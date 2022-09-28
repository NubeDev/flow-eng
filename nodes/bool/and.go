package bool

import (
	"github.com/NubeDev/flow-eng/node"
)

type And struct {
	*node.Spec
}

func NewAnd(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, and, category)
	if err != nil {
		return nil, err
	}
	return &And{body}, nil
}

func (inst *And) Process() {
	Process(inst)
}

func (inst *And) Cleanup() {}
