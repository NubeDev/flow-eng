package trigger

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
)

type Random struct {
	*node.Spec
	lastInput bool
}

func NewRandom(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, RandomFloat, Category)
	min := node.BuildInput(node.MinInput, node.TypeFloat, nil, body.Inputs, str.New("min"))
	max := node.BuildInput(node.MaxInput, node.TypeFloat, nil, body.Inputs, str.New("max"))
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, nil)
	inputs := node.BuildInputs(min, max, trigger)

	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp("When ‘trigger’ transitions from ‘false’ to ‘true’, a random number between ‘min’ and ‘max’ values is produced at ‘output’. The number of decimal places that ‘output’ values have can be set from settings.")

	node := &Random{body, true}
	node.SetSchema(node.buildSchema())
	return node, nil
}

func (inst *Random) Process() {
	min := inst.ReadPinOrSettingsFloat(node.MinInput)
	max := inst.ReadPinOrSettingsFloat(node.MaxInput)
	input, _ := inst.ReadPinAsBool(node.TriggerInput)

	if input && !inst.lastInput {
		settings, err := inst.getSettings(inst.GetSettings())
		if err != nil {
			log.Errorf("Random Node err: failed to get settings err:%s", err.Error())
			return
		}
		precision := settings.Precision
		inst.WritePinFloat(node.Outp, float.RandFloat(min, max), precision)
	}
	inst.lastInput = input
}

// Custom Node Settings Schema

type RandomSettingsSchema struct {
	Name      schemas.String `json:"name"`
	Precision schemas.Number `json:"precision"`
	Max       schemas.Number `json:"max"`
	Min       schemas.Number `json:"min"`
}

type RandomSettings struct {
	Name      string  `json:"name"`
	Precision int     `json:"precision"`
	Max       float64 `json:"max"`
	Min       float64 `json:"min"`
}

func (inst *Random) buildSchema() *schemas.Schema {
	props := &RandomSettingsSchema{}

	// name
	props.Name.Title = "Name"
	props.Name.Default = "Random"

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
		"ui:order": array.Slice{"name", "max", "min", "precision"},
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
