package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Divide struct {
	*node.Spec
}

func NewDivide(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, divide, category)
	if err != nil {
		return nil, err
	}
	return &Divide{body}, nil
}

func (inst *Divide) Process() {
	process(inst)
}
