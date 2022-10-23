package link

import (
	"github.com/NubeDev/flow-eng/node"
)

type Input struct {
	*node.Spec
}

func NewInput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkInput, category)
	topic := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	value := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(topic, value)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Input{body}, nil
}

func (inst *Input) Process() {
	in1 := inst.ReadPin(node.In)
	inst.WritePin(node.Out, in1)
	topic, _ := inst.ReadPinAsString(node.Topic)
	if topic != "" {
		getStore().Add(topic, in1)
	}
}
