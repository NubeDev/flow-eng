package trigger

import (
	"encoding/json"

	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
)

type Random struct {
	*node.Spec
	lastInput  bool
	lastOutput float64
}

func NewRandom(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, RandomFloat, Category)
	min := node.BuildInput(node.MinInput, node.TypeFloat, 0, body.Inputs, true, false)
	max := node.BuildInput(node.MaxInput, node.TypeFloat, 100, body.Inputs, true, false)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(min, max, trigger)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp("When ‘trigger’ transitions from ‘false’ to ‘true’, a random number between ‘min’ and ‘max’ values is produced at ‘output’. The number of decimal places that ‘output’ values have can be set from settings.")

	n := &Random{body, true, 0}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Random) Process() {
	min := inst.ReadPinOrSettingsFloat(node.MinInput)
	max := inst.ReadPinOrSettingsFloat(node.MaxInput)
	input, _ := inst.ReadPinAsBool(node.TriggerInput)

	if input && !inst.lastInput {
		settings, err := inst.getSettings(inst.GetSettings())
		if err != nil {
			log.Errorf("Random Node err: failed to get settings err:%s", err.Error())
			inst.WritePinNull(node.Out)
			return
		}
		precision := settings.Precision
		random := float.RandFloat(min, max)
		inst.WritePinFloat(node.Out, random, precision)
		inst.lastOutput = random
	}
	inst.lastInput = input
	inst.WritePinFloat(node.Out, inst.lastOutput)
}

// Custom Node Settings Schema

type RandomSettingsSchema struct {
	Precision schemas.Number `json:"precision"`
	Min       schemas.Number `json:"min"`
	Max       schemas.Number `json:"max"`
}

type RandomSettings struct {
	Precision int     `json:"precision"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
}

func (inst *Random) buildSchema() *schemas.Schema {
	props := &RandomSettingsSchema{}

	// decimals
	props.Precision.Title = "Precision / Decimal Places"
	props.Precision.Default = 2

	// range
	props.Max.Title = "Max"
	props.Max.Default = 100
	props.Min.Title = "Min"
	props.Min.Default = 0

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"min", "max", "precision"},
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

func (inst *Random) getSettings(body map[string]interface{}) (*RandomSettings, error) {
	settings := &RandomSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
