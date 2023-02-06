package latch

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
)

type SetResetLatch struct {
	*node.Spec
	currentVal bool
}

func NewSetResetLatch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, setResetLatch, category)
	set := node.BuildInput(node.Set, node.TypeBool, nil, body.Inputs, false)     // TODO: this input shouldn't have a manual override value
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false) // TODO: this input shouldn't have a manual override value

	inputs := node.BuildInputs(set, reset)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &SetResetLatch{body, false}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *SetResetLatch) Process() {
	set, _ := inst.ReadPinAsBool(node.Set)
	reset, _ := inst.ReadPinAsBool(node.Reset)

	settings, err := inst.getSettings(inst.GetSettings())
	if err != nil {
		log.Errorf("Set-Reset Latch err: failed to get settings err:%s", err.Error())
		return
	}
	allowResetOnSetTrue := settings.ResetWhenTrue

	if set && !reset {
		inst.currentVal = true
	} else if allowResetOnSetTrue && reset && inst.currentVal {
		inst.currentVal = false
	} else if !set && inst.currentVal && reset {
		inst.currentVal = false
	}
	inst.WritePinBool(node.Out, inst.currentVal)
}

// Custom Node Settings Schema

type SetResetLatchSettingsSchema struct {
	ResetWhenTrue schemas.Boolean `json:"reset_when_true"`
}

type SetResetLatchSettings struct {
	ResetWhenTrue bool `json:"reset_when_true"`
}

func (inst *SetResetLatch) buildSchema() *schemas.Schema {
	props := &SetResetLatchSettingsSchema{}

	// Step Size
	props.ResetWhenTrue.Title = "`reset` when `set` is `true`?`"
	props.ResetWhenTrue.Default = false

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"reset_when_true"},
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

func (inst *SetResetLatch) getSettings(body map[string]interface{}) (*SetResetLatchSettings, error) {
	settings := &SetResetLatchSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
