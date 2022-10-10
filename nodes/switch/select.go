package switches

import (
	"github.com/NubeDev/flow-eng/node"
)

type SelectNum struct {
	*node.Spec
}

func NewSelectNum(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, selectNum, category)
	if err != nil {
		return nil, err
	}
	return &SelectNum{body}, nil
}

func (inst *SelectNum) Process() {
	process(inst)
	//go getPoints()
}
