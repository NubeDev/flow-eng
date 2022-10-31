package filter

import (
	"github.com/NubeDev/flow-eng/node"
)

type FilterBool struct {
	*node.Spec
}

func NewFilterBool(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, sin, category)
	if err != nil {
		return nil, err
	}
	return &FilterBool{body}, nil
}

func (inst *FilterBool) Process() {
	process(inst)
}
