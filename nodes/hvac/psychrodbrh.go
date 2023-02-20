package hvac

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/psychrometrics"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/enescakir/emoji"
)

type PsychroDBRH struct {
	*node.Spec
	unitSystem   string
	isImperial   bool
	lastAltitude float64
}

func NewPsychroDBRH(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, psychroDBRH, category)
	dryBulbT := node.BuildInput(node.DryBulbTemp, node.TypeFloat, nil, body.Inputs, false, false)
	relHumPerc := node.BuildInput(node.RelHumid, node.TypeFloat, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(dryBulbT, relHumPerc)

	humRatio := node.BuildOutput(node.HumRatioO, node.TypeFloat, nil, body.Outputs)
	wetBulbT := node.BuildOutput(node.WetBulbTempO, node.TypeFloat, nil, body.Outputs)
	dewPointT := node.BuildOutput(node.DewPointTempO, node.TypeFloat, nil, body.Outputs)
	vapPres := node.BuildOutput(node.VaporPres, node.TypeFloat, nil, body.Outputs)
	enthalpy := node.BuildOutput(node.MoistAirEnthalpy, node.TypeFloat, nil, body.Outputs)
	volume := node.BuildOutput(node.MoistAirVolume, node.TypeFloat, nil, body.Outputs)
	degSaturation := node.BuildOutput(node.DegreeSaturation, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(humRatio, wetBulbT, dewPointT, vapPres, enthalpy, volume, degSaturation)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &PsychroDBRH{body, "Metric/SI", false, 1}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *PsychroDBRH) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())
	units := settings.UnitSystem
	altitude := settings.Altitude

	if units != inst.unitSystem || altitude != inst.lastAltitude {
		inst.setSubtitle(units, altitude)
		inst.unitSystem = units
		inst.lastAltitude = altitude

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
	dryBulbT, _ := inst.ReadPinAsFloat(node.DryBulbTemp)
	relHumPerc, _ := inst.ReadPinAsFloat(node.RelHumid)
	relHumPerc = relHumPerc / 100

	HumRatio, TWetBulb, TDewPoint, VapPres, MoistAirEnthalpy, MoistAirVolume, DegreeOfSaturation, err := psychrometrics.CalcPsychrometricsFromRelHum(dryBulbT, relHumPerc, atmPress, inst.isImperial)
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
	inst.WritePinFloat(node.DewPointTempO, TDewPoint, 4)
	inst.WritePinFloat(node.VaporPres, VapPres, 4)
	inst.WritePinFloat(node.MoistAirEnthalpy, MoistAirEnthalpy/1000, 4)
	inst.WritePinFloat(node.MoistAirVolume, MoistAirVolume, 4)
	inst.WritePinFloat(node.DegreeSaturation, DegreeOfSaturation, 4)

}

func (inst *PsychroDBRH) setSubtitle(units string, altitude float64) {
	subtitleText := ""
	if units == "Metric/SI" {
		subtitleText = fmt.Sprintf("units %v   altitude: %vm", units, altitude)
	} else {
		subtitleText = fmt.Sprintf("units %v   altitude: %vf", units, altitude)
	}
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type PsychroDBRHSettingsSchema struct {
	UnitSystem schemas.EnumString `json:"unit-system"`
	Altitude   schemas.Number     `json:"altitude"`
}

type PsychroDBRHSettings struct {
	UnitSystem string  `json:"unit-system"`
	Altitude   float64 `json:"altitude"`
}

func (inst *PsychroDBRH) buildSchema() *schemas.Schema {
	props := &PsychroDBRHSettingsSchema{}

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

func (inst *PsychroDBRH) getSettings(body map[string]interface{}) (*PsychroDBRHSettings, error) {
	settings := &PsychroDBRHSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
