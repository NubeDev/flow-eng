package hvac

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/pid"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type PIDNode struct {
	*node.Spec
	PID       *pid.Pid
	lastTime  int64
	startTime int64
	lastReset bool
}

func NewPIDNode(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pidNode, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, str.New("enable"))
	processValue := node.BuildInput(node.ProcessValue, node.TypeFloat, nil, body.Inputs, nil)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, nil, body.Inputs, str.New("setpoint"))
	minOut := node.BuildInput(node.MinOut, node.TypeFloat, nil, body.Inputs, str.New("min_out"))
	maxOut := node.BuildInput(node.MaxOut, node.TypeFloat, nil, body.Inputs, str.New("max_out"))
	inP := node.BuildInput(node.InP, node.TypeFloat, nil, body.Inputs, str.New("in_p"))
	inI := node.BuildInput(node.InI, node.TypeFloat, nil, body.Inputs, str.New("in_i"))
	inD := node.BuildInput(node.InD, node.TypeFloat, nil, body.Inputs, str.New("in_d"))
	direction := node.BuildInput(node.PIDDirection, node.TypeBool, nil, body.Inputs, str.New("direction"))
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, str.New("interval"))
	bias := node.BuildInput(node.Bias, node.TypeFloat, nil, body.Inputs, str.New("bias"))
	manual := node.BuildInput(node.Manual, node.TypeFloat, nil, body.Inputs, str.New("manual"))
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, nil)
	inputs := node.BuildInputs(enable, processValue, setPoint, minOut, maxOut, inP, inI, inD, direction, interval, bias, manual, reset)

	output := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(output)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &PIDNode{body, nil, 0, 0, false}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *PIDNode) Process() {
	if inst.PID == nil {
		inst.PID = pid.NewPid(0, 0, 1, 0, 0, 10, pid.DIRECT)
	}

	reset, _ := inst.ReadPinAsBool(node.Reset)
	if reset && !inst.lastReset {
		inst.PID.Initialize()
	}
	inst.lastReset = reset

	input, inputNull := inst.ReadPinAsFloat(node.ProcessValue)
	setpoint := inst.ReadPinOrSettingsFloat(node.Setpoint)
	enable := inst.ReadPinOrSettingsBool(node.Enable)

	if !enable || inputNull {
		inst.PID.SetMode(pid.MANUAL)
		manual := inst.ReadPinOrSettingsFloat(node.Manual)
		inst.WritePinFloat(node.Outp, manual)
		return
	}

	inst.PID.SetMode(pid.AUTO)
	inst.PID.SetSetpoint(setpoint)
	inst.PID.SetInput(input)

	minOut := inst.ReadPinOrSettingsFloat(node.MinOut)
	maxOut := inst.ReadPinOrSettingsFloat(node.MaxOut)
	inst.PID.SetOutputLimits(minOut, maxOut)

	inP := inst.ReadPinOrSettingsFloat(node.InP)
	inI := inst.ReadPinOrSettingsFloat(node.InI)
	inD := inst.ReadPinOrSettingsFloat(node.InD)
	inst.PID.SetTunings(inP, inI, inD)

	dir := inst.ReadPinOrSettingsBool(node.PIDDirection)
	inst.PID.SetControllerDirection(pid.PID_DIRECTION(dir))

	interval, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if interval <= 0 {
		interval = time.Second * 10
	}
	inst.PID.SetSampleTime(float64(interval.Milliseconds()))

	bias := inst.ReadPinOrSettingsFloat(node.Bias)
	inst.PID.SetBias(bias)

	inst.PID.Compute()
	inst.WritePinFloat(node.Outp, inst.PID.GetOutput())
}

func (inst *PIDNode) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type PIDNodeSettingsSchema struct {
	Enable            schemas.Boolean    `json:"enable"`
	Setpoint          schemas.Number     `json:"setpoint"`
	MinOut            schemas.Number     `json:"min_out"`
	MaxOut            schemas.Number     `json:"max_out"`
	InP               schemas.Number     `json:"in_p"`
	InI               schemas.Number     `json:"in_i"`
	InD               schemas.Number     `json:"in_d"`
	Direction         schemas.Boolean    `json:"direction"`
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	Bias              schemas.Number     `json:"bias"`
	Manual            schemas.Number     `json:"manual"`
}

type PIDNodeSettings struct {
	Enable            bool    `json:"enable"`
	Setpoint          float64 `json:"setpoint"`
	MinOut            float64 `json:"min_out"`
	MaxOut            float64 `json:"max_out"`
	InP               float64 `json:"in_p"`
	InI               float64 `json:"in_i"`
	InD               float64 `json:"in_d"`
	Direction         bool    `json:"direction"`
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
	Bias              float64 `json:"bias"`
	Manual            float64 `json:"manual"`
}

func (inst *PIDNode) buildSchema() *schemas.Schema {
	props := &PIDNodeSettingsSchema{}

	// enable
	props.Enable.Title = "Enable"
	props.Enable.Default = false

	// setpoint
	props.Setpoint.Title = "Setpoint"
	props.Setpoint.Default = 0

	// output limits
	props.MinOut.Title = "Minimum Output"
	props.MinOut.Default = 0
	props.MaxOut.Title = "Maximum Output"
	props.MaxOut.Default = 100

	// tuning factors
	props.InP.Title = "Proportional Factor (error multiplier)"
	props.InP.Default = 1
	props.InI.Title = "Integral Factor (repeats per second)"
	props.InI.Default = 0
	props.InD.Title = "Derivative Factor"
	props.InD.Default = 0

	// loop direction
	props.Direction.Title = "Direction (Direct/Reverse)"
	props.Direction.Default = true
	props.Direction.EnumNames = []string{"True: Reverse/Heating (increase when pv < sp)", "False: Direct/Cooling (increase when pv > sp)"}

	// time selection
	props.Interval.Title = "Update Interval"
	props.Interval.Default = 10
	props.IntervalTimeUnits.Title = "Update Interval Units"
	props.IntervalTimeUnits.Default = ttime.Sec
	props.IntervalTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// bias
	props.Bias.Title = "Bias (initial output)"
	props.Bias.Default = 0

	// manual
	props.Bias.Title = "Manual (output when disabled)"
	props.Bias.Default = 0

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"direction": array.Map{
			"ui:widget": "select",
		},
		"ui:order": array.Slice{"enable", "setpoint", "min_out", "max_out", "in_p", "in_i", "in_d", "direction", "interval", "interval_time_units", "bias", "manual"},
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

func (inst *PIDNode) getSettings(body map[string]interface{}) (*PIDNodeSettings, error) {
	settings := &PIDNodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
