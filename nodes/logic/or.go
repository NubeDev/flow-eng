package logic

import (
	"github.com/NubeDev/flow-eng/node"
)

type Or struct {
	*node.Spec
}

func NewOr(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, or, category)
	if err != nil {
		return nil, err
	}
	return &Or{body}, nil
}

func (inst *Or) Process() {
	Process(inst)
}

func (inst *Or) Cleanup() {}
