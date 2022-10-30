package link

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
)

type InputNum struct {
	*node.Spec
}

func NewInputNum(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkInputNum, category)
	topic := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	value := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(topic, value)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &InputNum{body}, nil
}

func (inst *InputNum) Process() {
	in1, _ := inst.ReadPinAsFloat(node.In)
	inst.WritePin(node.Out, in1)
	topic, _ := inst.ReadPinAsString(node.Topic)
	if topic != "" {
		topic = fmt.Sprintf("num-%s", topic)
		getStore().Add(topic, in1)
	}
}

type OutputNum struct {
	*node.Spec
}

func NewOutputNum(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkOutputNum, category)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, nil, outputs, body.Settings)
	return &OutputNum{body}, nil
}

func (inst *OutputNum) Process() {
	topic, _ := getSettings(inst.GetSettings())
	v, found := getStore().Get(topic)
	if found {
		inst.WritePinFloat(node.Out, conversions.GetFloat(v))
	} else {
		inst.WritePinNull(node.Out)
	}
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.SetSubTitle(fmt.Sprintf("%s topic: %s", emoji.Label, topic))
	}
}
func (inst *OutputNum) GetSchema() *schemas.Schema {
	return buildSchema("num")
}
