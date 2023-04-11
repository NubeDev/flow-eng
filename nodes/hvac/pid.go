package hvac

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/pid"
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
	enable := node.BuildInput(node.Enable, node.TypeBool, false, body.Inputs, false, false)
	processValue := node.BuildInput(node.ProcessValue, node.TypeFloat, nil, body.Inputs, false, false)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, 0, body.Inputs, true, false)
	minOut := node.BuildInput(node.MinOut, node.TypeFloat, 0, body.Inputs, true, false)
	maxOut := node.BuildInput(node.MaxOut, node.TypeFloat, 100, body.Inputs, true, false)
	inP := node.BuildInput(node.InP, node.TypeFloat, 1, body.Inputs, true, false)
	inI := node.BuildInput(node.InI, node.TypeFloat, 0.01, body.Inputs, true, false)
	inD := node.BuildInput(node.InD, node.TypeFloat, 0, body.Inputs, true, false)
	direction := node.BuildInput(node.PIDDirection, node.TypeBool, false, body.Inputs, true, false)
	interval := node.BuildInput(node.Interval, node.TypeFloat, 10, body.Inputs, true, false)
	bias := node.BuildInput(node.Bias, node.TypeFloat, 0, body.Inputs, true, false)
	manual := node.BuildInput(node.Manual, node.TypeFloat, 0, body.Inputs, true, false)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(enable, processValue, setPoint, minOut, maxOut, inP, inI, inD, direction, interval, bias, manual, reset)

	output := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
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
	enable, _ := inst.ReadPinAsBool(node.Enable)

	minOut := inst.ReadPinOrSettingsFloat(node.MinOut)
	maxOut := inst.ReadPinOrSettingsFloat(node.MaxOut)

	if !enable || inputNull {
		inst.PID.SetMode(pid.MANUAL)
		manual := inst.ReadPinOrSettingsFloat(node.Manual)
		inst.WritePinFloat(node.Out, manual)
		return
	} else if minOut == maxOut {
		inst.PID.SetMode(pid.MANUAL)
		inst.WritePinFloat(node.Out, minOut)
		return
	}

	inst.PID.SetMode(pid.AUTO)
	inst.PID.SetSetpoint(setpoint)
	inst.PID.SetInput(input)
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
	inst.WritePinFloat(node.Out, inst.PID.GetOutput())
}

func (inst *PIDNode) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type PIDNodeSettingsSchema struct {
	Setpoint          schemas.Number     `json:"setpoint"`
	MinOut            schemas.Number     `json:"min-out"`
	MaxOut            schemas.Number     `json:"max-out"`
	InP               schemas.Number     `json:"in-p"`
	InI               schemas.Number     `json:"in-i"`
	InD               schemas.Number     `json:"in-d"`
	Direction         schemas.Boolean    `json:"direction"`
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	Bias              schemas.Number     `json:"bias"`
	Manual            schemas.Number     `json:"manual"`
}

type PIDNodeSettings struct {
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
	props.InI.Default = 0.01
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
		"ui:order": array.Slice{"enable", "setpoint", "min-out", "max-out", "in-p", "in-i", "in-d", "direction", "interval", "interval_time_units", "bias", "manual"},
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
