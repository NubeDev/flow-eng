package link

import (
	"github.com/NubeDev/flow-eng/node"
)

type Output struct {
	*node.Spec
}

func NewOutput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkOutput, category)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, nil, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &Output{body}, nil
}

func (inst *Output) Process() {
	topic, _ := getSettings(inst.GetSettings())
	v, found := getStore().Get(topic)
	if found {
		inst.WritePin(node.Out, v)
	} else {
		inst.WritePinNull(node.Out)
	}
}
