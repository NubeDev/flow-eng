package link

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
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
		topic = fmt.Sprintf("string-%s", topic)
		getStore().Add(topic, in1)
	}
}

type Output struct {
	*node.Spec
}

func NewOutput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkOutput, category)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, nil, outputs, body.Settings)
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
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.SetSubTitle(fmt.Sprintf("%s topic: %s", emoji.Label, topic))
	}
}

func (inst *Output) GetSchema() *schemas.Schema {
	return buildSchema("string")
}
