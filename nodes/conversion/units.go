package conversion

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/NubeIO/lib-units/units"
)

type Units struct {
	*node.Spec
	from  string
	to    string
	value string
}

const unitsDesc = `node used for converting values between metric and imperial and vice-versa`

func NewUnits(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, conversionUnit, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false))
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outErr := node.BuildOutput(node.ErrMsg, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(out, outErr)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	n := &Units{body, "", "", ""}
	n.SetSchema(n.buildSchema())
	n.SetHelp(unitsDesc)
	return n, nil
}

func (inst *Units) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.init()
	}
	in1, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		process, err := units.Process(in1, inst.from, inst.to)
		if err != nil {
			inst.WritePinNull(node.Out)
			inst.WritePinAsString(node.ErrMsg, err.Error())
			return
		}
		inst.value = process.String()
		inst.WritePinFloat(node.Out, process.AsFloat(), 2)
		inst.setSubtitle()
	}
}

func (inst *Units) setSubtitle() {
	title := fmt.Sprintf("%s to %s", inst.from, inst.to)
	inst.SetSubTitle(title)
}

func (inst *Units) init() {
	settings, _ := inst.getSettings(inst.GetSettings())
	if settings != nil {
		inst.from = settings.UnitFrom
		inst.to = settings.UnitTo
	}
}

type enumString struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:""`
	Help     string   `json:"help" default:""`
	Default  string   `json:"default" default:""`
	Options  []string `json:"enum" default:"[]"`
	EnumName []string `json:"enumNames" default:"[]"`
}

type unitsSettingsSchema struct {
	UnitFrom enumString `json:"unit_from"`
	UnitTo   enumString `json:"unit_to"`
}

type unitsSettings struct {
	UnitFrom string `json:"unit_from"`
	UnitTo   string `json:"unit_to"`
}

func (inst *Units) buildSchema() *schemas.Schema {
	props := &unitsSettingsSchema{}

	names, unit := units.SupportedUnitsNames()

	props.UnitFrom.Title = "Units From"
	props.UnitFrom.Options = unit
	props.UnitFrom.EnumName = names
	props.UnitFrom.Default = "disabled"

	props.UnitTo.Title = "Units To"
	props.UnitTo.Options = unit
	props.UnitTo.EnumName = names
	props.UnitTo.Default = "disabled"

	schema.Set(props)

	uiSchema := array.Map{
		"time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"unit_from", "unit_to"},
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
func (inst *Units) getSettings(body map[string]interface{}) (*unitsSettings, error) {
	settings := &unitsSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
