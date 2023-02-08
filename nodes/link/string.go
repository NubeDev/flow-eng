package link

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type StringLinkInput struct {
	*node.Spec
	lastTopic string
}

func NewStringLinkInput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkInputString, category)
	topic := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs, true)
	value := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false)
	inputs := node.BuildInputs(topic, value)
	body = node.BuildNode(body, inputs, nil, body.Settings)
	n := &StringLinkInput{body, ""}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *StringLinkInput) Process() {
	in1, _ := inst.ReadPinAsString(node.In)
	topic := inst.ReadPinOrSettingsString(node.Topic)
	if topic != inst.lastTopic {
		parentTopic := helpers.CleanParentName(topic, inst.GetParentName())
		if parentTopic != "" {
			topic = parentTopic
		}
		topic = fmt.Sprintf("string-%s", topic)
		getStore().Add(topic, in1)
		inst.SetSubTitle(topic)
		inst.lastTopic = topic
	}
}

// Custom Node Settings Schema

type StringLinkInputSettingsSchema struct {
	Topic schemas.String `json:"topic"`
}

type StringLinkInputSettings struct {
	Topic string `json:"topic"`
}

func (inst *StringLinkInput) buildSchema() *schemas.Schema {
	props := &StringLinkInputSettingsSchema{}

	// topic
	props.Topic.Title = "Topic"
	props.Topic.Default = ""

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"topic"},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Node Settings",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

func (inst *StringLinkInput) getSettings(body map[string]interface{}) (*StringLinkInputSettings, error) {
	settings := &StringLinkInputSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

type StringLinkOutput struct {
	*node.Spec
}

func NewStringLinkOutput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkOutputString, category)
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, nil, outputs, body.Settings)
	return &StringLinkOutput{body}, nil
}

func (inst *StringLinkOutput) Process() {
	topic, _ := getSettings(inst.GetSettings())
	v, found := getStore().Get(topic)
	if found {
		inst.WritePin(node.Out, v)
	} else {
		inst.WritePinNull(node.Out)
	}
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.SetSubTitle(fmt.Sprintf("topic: %s", topic))
	}
}

func (inst *StringLinkOutput) GetSchema() *schemas.Schema {
	return buildSchema("string")
}
