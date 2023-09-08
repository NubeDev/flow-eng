package math

import (
	"fmt"

	"github.com/NubeDev/flow-eng/node"
)

type Add struct {
	*node.Spec
}

func NewAdd(body *node.Spec, _ ...any) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, add, Category)
	if err != nil {
		return nil, err
	}
	body.SetHelp(fmt.Sprintf("%s addition", mathHelp))
	return &Add{body}, nil
}

func (inst *Add) Process() {
	process(inst)
}
