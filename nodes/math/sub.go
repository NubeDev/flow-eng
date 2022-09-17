package math

import (
	"github.com/NubeDev/flow-eng/node"
)

type Sub struct {
	*node.Spec
}

func NewSub(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, sub, category)
	if err != nil {
		return nil, err
	}
	return &Sub{body}, nil
}

func (inst *Sub) Process() {
	process(inst)
}

func (inst *Sub) Cleanup() {}
