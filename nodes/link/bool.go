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

type BoolLinkInput struct {
	*node.Spec
	lastTopic string
}

func NewBoolLinkInput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkInputBool, category)
	topic := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs, true, false)
	value := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(topic, value)
	body = node.BuildNode(body, inputs, nil, body.Settings)
	n := &BoolLinkInput{body, ""}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *BoolLinkInput) Process() {
	in1, _ := inst.ReadPinAsBool(node.In)
	topic := inst.ReadPinOrSettingsString(node.Topic)
	if topic != inst.lastTopic {
		parentTopic := helpers.CleanParentName(topic, inst.GetParentName())
		if parentTopic != "" {
			topic = parentTopic
		}
		topic = fmt.Sprintf("bool-%s", topic)
		getStore().Add(topic, in1)
		inst.SetSubTitle(topic)
		inst.lastTopic = topic
	}
}

// Custom Node Settings Schema

type BoolLinkInputSettingsSchema struct {
	Topic schemas.String `json:"topic"`
}

type BoolLinkInputSettings struct {
	Topic string `json:"topic"`
}

func (inst *BoolLinkInput) buildSchema() *schemas.Schema {
	props := &BoolLinkInputSettingsSchema{}

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

func (inst *BoolLinkInput) getSettings(body map[string]interface{}) (*BoolLinkInputSettings, error) {
	settings := &BoolLinkInputSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

type BoolLinkOutput struct {
	*node.Spec
}

func NewBoolLinkOutput(body *node.Spec, store *Store) (node.Node, error) {
	if store == nil {
		store = getStore()
	}
	body = node.Defaults(body, linkOutputBool, category)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, nil, outputs, body.Settings)
	return &BoolLinkOutput{body}, nil
}

func (inst *BoolLinkOutput) Process() {
	topic, _ := getSettings(inst.GetSettings())
	v, found := getStore().Get(topic)
	if found {
		inst.WritePinFloat(node.Out, conversions.GetFloat(v))
	} else {
		inst.WritePinNull(node.Out)
	}
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.SetSubTitle(fmt.Sprintf("topic: %s", topic))
	}
}
func (inst *BoolLinkOutput) GetSchema() *schemas.Schema {
	return buildSchema("bool")
}
