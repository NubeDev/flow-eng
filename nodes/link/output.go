package link

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
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
	topic := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(topic)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Output{body}, nil
}

func (inst *Output) Process() {
	topic := conversions.ToString(inst.ReadPin(node.Topic))
	v, _ := getStore().Get(topic)
	inst.WritePin(node.Out, v)
}
