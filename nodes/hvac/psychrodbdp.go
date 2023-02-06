package hvac

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/psychrometrics"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/enescakir/emoji"
)

type PsychroDBDP struct {
	*node.Spec
	unitSystem string
	isImperial bool
}

func NewPsychroDBDP(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, psychroDBDP, category)
	dryBulbT := node.BuildInput(node.DryBulbTemp, node.TypeFloat, nil, body.Inputs, nil)
	dewPointT := node.BuildInput(node.DewPointTemp, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(dryBulbT, dewPointT)

	humRatio := node.BuildOutput(node.HumRatioO, node.TypeFloat, nil, body.Outputs)
	wetBulbT := node.BuildOutput(node.WetBulbTempO, node.TypeFloat, nil, body.Outputs)
	relHumPerc := node.BuildOutput(node.RelHumPercO, node.TypeFloat, nil, body.Outputs)
	vapPres := node.BuildOutput(node.VaporPres, node.TypeFloat, nil, body.Outputs)
	enthalpy := node.BuildOutput(node.MoistAirEnthalpy, node.TypeFloat, nil, body.Outputs)
	volume := node.BuildOutput(node.MoistAirVolume, node.TypeFloat, nil, body.Outputs)
	degSaturation := node.BuildOutput(node.DegreeSaturation, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(humRatio, wetBulbT, relHumPerc, vapPres, enthalpy, volume, degSaturation)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &PsychroDBDP{body, "Metric/SI", false}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *PsychroDBDP) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())
	units := settings.UnitSystem
	altitude := settings.Altitude

	if units != inst.unitSystem {
		inst.setSubtitle()
		inst.unitSystem = units

	}
	if units == "Metric/SI" {
		inst.isImperial = false
	} else {
		inst.isImperial = true
	}
	atmPress, err := psychrometrics.GetStandardAtmPressure(altitude, inst.isImperial)
	if err != nil {
		inst.SetWaringMessage(err.Error())
		inst.SetWaringIcon(string(emoji.RedCircle))
		return
	}
	dryBulbT := inst.ReadPinOrSettingsFloat(node.DryBulbTemp)
	dewPointT := inst.ReadPinOrSettingsFloat(node.DewPointTemp)

	HumRatio, TWetBulb, RelHum, VapPres, MoistAirEnthalpy, MoistAirVolume, DegreeOfSaturation, err := psychrometrics.CalcPsychrometricsFromTDewPoint(dryBulbT, dewPointT, atmPress, inst.isImperial)
	if err != nil {
		inst.SetWaringMessage(err.Error())
		inst.SetWaringIcon(string(emoji.RedCircle))
		return
	} else {
		inst.SetWaringMessage("")
		inst.SetWaringIcon(string(emoji.GreenCircle))
	}

	inst.WritePinFloat(node.HumRatioO, HumRatio, 4)
	inst.WritePinFloat(node.WetBulbTempO, TWetBulb, 4)
	inst.WritePinFloat(node.RelHumPercO, RelHum*100, 4)
	inst.WritePinFloat(node.VaporPres, VapPres, 4)
	inst.WritePinFloat(node.MoistAirEnthalpy, MoistAirEnthalpy/1000, 4)
	inst.WritePinFloat(node.MoistAirVolume, MoistAirVolume, 4)
	inst.WritePinFloat(node.DegreeSaturation, DegreeOfSaturation, 4)

}

func (inst *PsychroDBDP) setSubtitle() {
	subtitleText := inst.unitSystem
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type PsychroDBDPSettingsSchema struct {
	UnitSystem schemas.EnumString `json:"unit-system"`
	Altitude   schemas.Number     `json:"altitude"`
}

type PsychroDBDPSettings struct {
	UnitSystem string  `json:"unit-system"`
	Altitude   float64 `json:"altitude"`
}

func (inst *PsychroDBDP) buildSchema() *schemas.Schema {
	props := &PsychroDBDPSettingsSchema{}

	// unit system
	props.UnitSystem.Title = "Unit System"
	props.UnitSystem.Default = "Metric/SI"
	props.UnitSystem.Options = []string{"Metric/SI", "Imperial/IP"}
	props.UnitSystem.EnumName = []string{"Metric/SI", "Imperial/IP"}

	// altitude
	props.UnitSystem.Title = "Altitude (m/ft)"

	schema.Set(props)

	uiSchema := array.Map{
		"unit-system": array.Map{
			"ui:widget": "select",
		},
		"ui:order": array.Slice{"unit-system", "altitude"},
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

func (inst *PsychroDBDP) getSettings(body map[string]interface{}) (*PsychroDBDPSettings, error) {
	settings := &PsychroDBDPSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
