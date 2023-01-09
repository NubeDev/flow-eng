package statistics

import (
	"github.com/NubeDev/flow-eng/node"
)

type Avg struct {
	*node.Spec
}

func NewAvg(body *node.Spec) (node.Node, error) {
	var err error
	body, err = nodeDefault(body, avg, category)
	if err != nil {
		return nil, err
	}
	return &Avg{body}, nil
}

func (inst *Avg) Process() {
	process(inst)
}
