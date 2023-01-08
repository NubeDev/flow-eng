package math

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

type Add struct {
	*node.Spec
}

func NewAdd(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, add, category)
	if err != nil {
		return nil, err
	}
	body.SetHelp(fmt.Sprintf("Add %s", mathHelp))
	return &Add{body}, nil
}

func (inst *Add) Process() {
	process(inst)
}
