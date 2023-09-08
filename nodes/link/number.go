package link

import (
	"encoding/json"
	"fmt"

	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type NumLinkInput struct {
	*node.Spec
	lastTopic string
}

func NewNumLinkInput(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, linkInputNum, Category)
	topic := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs, true, false)
	value := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(topic, value)
	body = node.BuildNode(body, inputs, nil, body.Settings)

	n := &NumLinkInput{body, ""}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *NumLinkInput) Process() {
	in1, null := inst.ReadPinAsFloat(node.In)
	topic := inst.ReadPinOrSettingsString(node.Topic)
	if topic != inst.lastTopic {
		parentTopic := helpers.CleanParentName(topic, inst.GetParentName())
		if parentTopic != "" {
			topic = parentTopic
		}
		topic = fmt.Sprintf("num-%s", topic)
		if null {
			getStore().Add(topic, nil)
		} else {
			getStore().Add(topic, in1)
		}

		inst.SetSubTitle(topic)
		inst.lastTopic = topic
	}
}

// Custom Node Settings Schema

type NumLinkInputSettingsSchema struct {
	Topic schemas.String `json:"topic"`
}

type NumLinkInputSettings struct {
	Topic string `json:"topic"`
}

func (inst *NumLinkInput) buildSchema() *schemas.Schema {
	props := &NumLinkInputSettingsSchema{}

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

func (inst *NumLinkInput) getSettings(body map[string]interface{}) (*NumLinkInputSettings, error) {
	settings := &NumLinkInputSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

type NumLinkOutput struct {
	*node.Spec
}

func NewNumLinkOutput(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, linkOutputNum, Category)
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, nil, outputs, body.Settings)
	return &NumLinkOutput{body}, nil
}

func (inst *NumLinkOutput) Process() {
	topic, _ := getSettings(inst.GetSettings())
	v, found := getStore().Get(topic)
	if found {
		if v == nil {
			inst.WritePinNull(node.Out)
		} else {
			inst.WritePinFloat(node.Out, conversions.GetFloat(v))
		}
	} else {
		inst.WritePinNull(node.Out)
	}
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.SetSubTitle(fmt.Sprintf("topic: %s", topic))
	}
}
func (inst *NumLinkOutput) GetSchema() *schemas.Schema {
	return buildSchema("num")
}
