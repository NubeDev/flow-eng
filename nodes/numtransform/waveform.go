package numtransform

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type Waveform struct {
	*node.Spec
	startTime      time.Time
	lastTime       time.Time
	lastInterval   time.Duration
	lastPeriod     time.Duration
	lastEnable     bool
	lastOutput     float64
	lastSignalType string
	lastInvert     bool
	lastAmplitude  float64
	lastReset      bool
}

func NewWaveform(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, waveform, Category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, false)
	interval := node.BuildInput(node.Interval, node.TypeFloat, 1, body.Inputs, true, false)
	period := node.BuildInput(node.Period, node.TypeFloat, 60, body.Inputs, true, false)
	amplitude := node.BuildInput(node.Amplitude, node.TypeFloat, 100, body.Inputs, true, false)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false, true)
	inputs := node.BuildInputs(enable, interval, period, amplitude, reset)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &Waveform{body, time.Now(), time.Now(), 5 * time.Second, 1 * time.Second, true, 0, "", true, 0, true}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Waveform) Process() {
	enable, _ := inst.ReadPinAsBool(node.Enable)
	reset, _ := inst.ReadPinAsBool(node.Reset)
	if !enable {
		inst.WritePinFloat(node.Out, inst.lastOutput)
	} else if enable && !inst.lastEnable || reset && !inst.lastReset {
		inst.startTime = time.Now()
	}
	inst.lastEnable = enable
	inst.lastReset = reset

	settings, _ := inst.getSettings(inst.GetSettings())
	signalType := settings.SignalType
	invert := settings.Invert
	invertVal := 1
	if invert {
		invertVal = -1
	}

	amplitude := inst.ReadPinOrSettingsFloat(node.Amplitude)
	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	periodDuration, _ := inst.ReadPinAsTimeSettings(node.Period)
	if intervalDuration != inst.lastInterval || periodDuration != inst.lastPeriod || signalType != inst.lastSignalType || invert != inst.lastInvert || amplitude != inst.lastAmplitude {
		inst.setSubtitle(intervalDuration, periodDuration, signalType, invert, amplitude)
		inst.lastInterval = intervalDuration
		inst.lastPeriod = periodDuration
		inst.lastSignalType = signalType
		inst.lastInvert = invert
		inst.lastAmplitude = amplitude
	}

	output := inst.lastOutput
	if inst.lastTime.Add(intervalDuration).Before(time.Now()) {
		waveformProgress := time.Now().Unix() - inst.startTime.Unix()
		t := math.Mod(float64(waveformProgress), periodDuration.Seconds()) / periodDuration.Seconds()
		switch signalType {
		case "Ramp":
			output = 2 * math.Abs(math.Round(t)-t)

		case "Sine":
			output = math.Sin(2 * math.Pi * t)

		case "Square":
			sign := math.Signbit(math.Sin(2 * math.Pi * t))
			if sign {
				output = 0
			} else {
				output = 1
			}

		case "Triangle":
			output = 1 - 4*math.Abs(math.Round(t-0.25)-(t-0.25))

		case "Sawtooth":
			output = t
		}
		output = output * float64(invertVal) * amplitude
		inst.lastTime = time.Now()
	}
	inst.WritePinFloat(node.Out, output)
	inst.lastOutput = output
}

func (inst *Waveform) setSubtitle(intervalDuration, periodDuration time.Duration, signal string, invert bool, amplitude float64) {
	subtitleText := fmt.Sprintf("waveform: %s, amplitude: %v, update interval: %s, period: %s, invert: %t", signal, amplitude, intervalDuration.String(), periodDuration.String(), invert)
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type WaveformSettingsSchema struct {
	SignalType        schemas.EnumString `json:"waveform-type"`
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	Period            schemas.Number     `json:"period"`
	PeriodTimeUnits   schemas.EnumString `json:"period_time_units"`
	Amplitude         schemas.Number     `json:"amplitude"`
	Invert            schemas.Boolean    `json:"invert"`
}

type WaveformSettings struct {
	SignalType        string  `json:"waveform-type"`
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
	Period            float64 `json:"period"`
	PeriodTimeUnits   string  `json:"period_time_units"`
	Amplitude         float64 `json:"amplitude"`
	Invert            bool    `json:"invert"`
}

func (inst *Waveform) buildSchema() *schemas.Schema {
	props := &WaveformSettingsSchema{}

	// waveform type
	props.SignalType.Title = "Waveform Type"
	props.SignalType.Default = "Ramp"
	props.SignalType.Options = []string{"Ramp", "Sine", "Square", "Triangle", "Sawtooth"}
	props.SignalType.EnumName = []string{"Ramp", "Sine", "Square", "Triangle", "Sawtooth"}

	// interval
	props.Interval.Title = "Interval"
	props.Interval.Default = 1

	// interval time selection
	props.IntervalTimeUnits.Title = "Interval Units"
	props.IntervalTimeUnits.Default = ttime.Sec
	props.IntervalTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// period
	props.Period.Title = "Period"
	props.Period.Default = 30

	// period time selection
	props.PeriodTimeUnits.Title = "Period Units"
	props.PeriodTimeUnits.Default = ttime.Sec
	props.PeriodTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.PeriodTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// amplitude
	props.Amplitude.Title = "Amplitude"
	props.Amplitude.Default = 1

	// invert
	props.Invert.Title = "Invert Waveform"
	props.Invert.Default = false

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"period_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"waveform-type", "amplitude", "interval", "interval_time_units", "period", "period_time_units", "invert"},
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

func (inst *Waveform) getSettings(body map[string]interface{}) (*WaveformSettings, error) {
	settings := &WaveformSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
