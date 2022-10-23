package link

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
)

type Output struct {
	*node.Spec
}

func NewOutput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkOutput, category)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	topic := node.BuildOutput(node.OutTopic, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out, topic)
	body = node.BuildNode(body, nil, outputs, body.Settings)
	return &Output{body}, nil
}

func (inst *Output) Process() {
	topic, _ := getSettings(inst.GetSettings())
	v, found := getStore().Get(topic)
	if found {
		inst.WritePin(node.Out, v)
		inst.WritePin(node.OutTopic, topic)
	} else {
		inst.WritePinNull(node.Out)
	}
}

func (inst *Output) GetSchema() *schemas.Schema {
	return inst.buildSchema()
}
