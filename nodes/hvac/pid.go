package hvac

import (
	"github.com/NubeDev/flow-eng/helpers/pid"
	"github.com/NubeDev/flow-eng/node"
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
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, nil)
	processValue := node.BuildInput(node.ProcessValue, node.TypeFloat, nil, body.Inputs, nil)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, nil, body.Inputs, nil)
	minOut := node.BuildInput(node.MinOut, node.TypeFloat, nil, body.Inputs, nil)
	maxOut := node.BuildInput(node.MaxOut, node.TypeFloat, nil, body.Inputs, nil)
	inP := node.BuildInput(node.InP, node.TypeFloat, nil, body.Inputs, nil)
	inI := node.BuildInput(node.InI, node.TypeFloat, nil, body.Inputs, nil)
	inD := node.BuildInput(node.InD, node.TypeFloat, nil, body.Inputs, nil)
	direction := node.BuildInput(node.PIDDirection, node.TypeBool, nil, body.Inputs, nil)
	intervalSecs := node.BuildInput(node.IntervalSecs, node.TypeFloat, nil, body.Inputs, nil)
	bias := node.BuildInput(node.Bias, node.TypeFloat, nil, body.Inputs, nil)
	manual := node.BuildInput(node.Manual, node.TypeFloat, nil, body.Inputs, nil)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs, nil)

	output := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)

	inputs := node.BuildInputs(enable, processValue, setPoint, minOut, maxOut, inP, inI, inD, direction, intervalSecs, bias, manual, reset)
	outputs := node.BuildOutputs(output)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	// body.SetSchema(buildSchema())
	return &PIDNode{body, nil, 0, 0, false}, nil
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
	setpoint, setpointNull := inst.ReadPinAsFloat(node.Setpoint)
	enable, _ := inst.ReadPinAsBool(node.Enable)

	if !enable || inputNull || setpointNull {
		inst.PID.SetMode(pid.MANUAL)
		manual, _ := inst.ReadPinAsFloat(node.Manual)
		inst.WritePinFloat(node.Out, manual)
		return
	}

	inst.PID.SetMode(pid.AUTO)
	inst.PID.SetSetpoint(setpoint)
	inst.PID.SetInput(input)

	minOut, null := inst.ReadPinAsFloat(node.MinOut)
	if null {
		minOut = 0
	}
	maxOut, null := inst.ReadPinAsFloat(node.MaxOut)
	if null {
		minOut = 100
	}
	inst.PID.SetOutputLimits(minOut, maxOut)

	inP, null := inst.ReadPinAsFloat(node.InP)
	if null {
		inP = 1
	}
	inI, null := inst.ReadPinAsFloat(node.InI)
	if null {
		inI = 0
	}
	inD, null := inst.ReadPinAsFloat(node.InD)
	if null {
		inD = 0
	}
	inst.PID.SetTunings(inP, inI, inD)

	direction := pid.DIRECT
	dir, _ := inst.ReadPinAsBool(node.PIDDirection)
	if dir {
		direction = pid.REVERSE
	}
	inst.PID.SetControllerDirection(direction)

	intervalSecs, _ := inst.ReadPinAsFloat(node.IntervalSecs)
	if intervalSecs <= 0 {
		intervalSecs = 10
	} else if intervalSecs > 500 {
		intervalSecs = 500
	}
	intervalMillis := intervalSecs * 1000
	inst.PID.SetSampleTime(intervalMillis)

	bias, _ := inst.ReadPinAsFloat(node.Bias)
	inst.PID.SetBias(bias)

	inst.PID.Compute()
	inst.WritePinFloat(node.Out, inst.PID.GetOutput())
}
