package boolean

import (
	"github.com/NubeDev/flow-eng/node"
)

type Or struct {
	*node.Spec
}

func NewOr(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, or, Category)
	if err != nil {
		return nil, err
	}
	return &Or{body}, nil
}

func (inst *Or) Process() {
	Process(inst)
}
