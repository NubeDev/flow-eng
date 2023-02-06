package numtransform

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type Fade struct {
	*node.Spec
	outEqTo      bool
	fading       bool
	originalFrom float64
	targetTo     float64
	startTime    time.Time
	lastInterval time.Duration
	fadeInterval time.Duration
}

func NewFade(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, fade, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, nil)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, str.New("interval"))
	from := node.BuildInput(node.From, node.TypeFloat, nil, body.Inputs, nil)
	to := node.BuildInput(node.To, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(enable, interval, from, to)

	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	outEqTo := node.BuildOutput(node.OutEqTo, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out, outEqTo)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &Fade{body, false, false, 0, 0, time.Time{}, 1 * time.Second, 1 * time.Second}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Fade) Process() {
	enable, _ := inst.ReadPinAsBool(node.Enable)
	from, fromNull := inst.ReadPinAsFloat(node.From)
	to, toNull := inst.ReadPinAsFloat(node.To)
	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if intervalDuration != inst.lastInterval {
		inst.setSubtitle(intervalDuration)
		inst.lastInterval = intervalDuration
	}
	if enable && !inst.fading && !inst.outEqTo && !toNull && !fromNull {
		inst.WritePinBool(node.OutEqTo, false)
		inst.WritePinFloat(node.Outp, from)
		inst.fading = true
		inst.startTime = time.Now()
		inst.targetTo = to
		inst.originalFrom = from
		inst.fadeInterval = intervalDuration
	} else if !enable || toNull {
		inst.WritePinBool(node.OutEqTo, false)
		inst.WritePinFloat(node.Outp, 0)
		inst.fading = false
	} else if enable && inst.fading {
		now := time.Now()
		if now.After(inst.startTime.Add(intervalDuration)) {
			inst.fading = false
			inst.outEqTo = true
			inst.WritePinBool(node.OutEqTo, true)
			inst.WritePinFloat(node.Outp, to)
		} else {
			inst.WritePinBool(node.OutEqTo, false)
			percThroughFade := 1 - float64(inst.startTime.Add(inst.fadeInterval).Unix()-now.Unix())/inst.fadeInterval.Seconds()
			fadingOutput := ((inst.targetTo - inst.originalFrom) * percThroughFade) + inst.originalFrom
			inst.WritePinBool(node.OutEqTo, false)
			inst.WritePinFloat(node.Outp, fadingOutput)
		}
	} else if enable && inst.outEqTo {
		inst.WritePinBool(node.OutEqTo, true)
		inst.WritePinFloat(node.Outp, to)
	}
}

func (inst *Fade) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type FadeSettingsSchema struct {
	Interval  schemas.Number     `json:"interval"`
	TimeUnits schemas.EnumString `json:"interval_time_units"`
}

type FadeSettings struct {
	Interval  float64 `json:"interval"`
	TimeUnits string  `json:"interval_time_units"`
}

func (inst *Fade) buildSchema() *schemas.Schema {
	props := &FadeSettingsSchema{}
	// time selection
	props.Interval.Title = "Interval"
	props.Interval.Default = 30

	// time selection
	props.TimeUnits.Title = "Interval Units"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.TimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "interval_time_units"},
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

func (inst *Fade) getSettings(body map[string]interface{}) (*FadeSettings, error) {
	settings := &FadeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
