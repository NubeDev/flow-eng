package timing

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type DutyCycle struct {
	*node.Spec
	onTicker         *time.Ticker
	offTicker        *time.Ticker
	offDelayTimer    *time.Timer
	cancelChannel    chan bool
	lastIntervalSecs float64
	lastDuty         float64
	lastEnable       bool
}

func NewDutyCycle(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, dutyCycle, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)
	intervalSecs := node.BuildInput(node.IntervalSecs, node.TypeFloat, nil, body.Inputs)
	dutyCycleInput := node.BuildInput(node.DutyCycle, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(enable, intervalSecs, dutyCycleInput)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &DutyCycle{body, nil, nil, nil, nil, 10, 50, false}, nil
}

/*
if node input is true
start delay, after the delay set the triggered to true
*/

func (inst *DutyCycle) Process() {
	// settings, _ := getSettings(inst.GetSettings())

	enable, _ := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.disableDutyCycle()
		return
	}

	intervalSecs, _ := inst.ReadPinAsFloat(node.Interval)
	if intervalSecs <= 0 {
		intervalSecs = 10
	}

	dutyCyclePerc, dutyNull := inst.ReadPinAsFloat(node.DutyCycle)
	if dutyCyclePerc < 0 || dutyCyclePerc > 100 || dutyNull {
		dutyCyclePerc = 50
	}

	// Check if there are settings that require a restart
	if enable && !inst.lastEnable || intervalSecs != inst.lastIntervalSecs || dutyCyclePerc != inst.lastDuty {
		inst.restartDutyCycle(intervalSecs, dutyCyclePerc)
	}
	inst.lastIntervalSecs = intervalSecs
	inst.lastDuty = dutyCyclePerc
	inst.lastEnable = enable

}

func (inst *DutyCycle) restartDutyCycle(intervalSeconds, dutyCycle float64) error {
	if intervalSeconds <= 0 || (dutyCycle < 0 || dutyCycle > 100) {
		return errors.New("restartDutyCycle() err: invalid inputs")
	}
	intervalDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", intervalSeconds))

	delayBetweenOnAndOff := intervalSeconds * (dutyCycle / 100)
	delayBetweenOnAndOffDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", delayBetweenOnAndOff))

	cancel := make(chan bool)
	inst.cancelChannel = cancel
	inst.onTicker = time.NewTicker(intervalDuration)
	inst.offDelayTimer = time.AfterFunc(delayBetweenOnAndOffDuration, func() {
		inst.offTicker = time.NewTicker(intervalDuration)
	})

	go func() {
		for {
			select {
			case <-cancel:
				return
			case <-inst.onTicker.C:
				inst.WritePinTrue(node.Out)
			case <-inst.offTicker.C:
				inst.WritePinFalse(node.Out)
			}
		}
	}()
	return nil
}

func (inst *DutyCycle) disableDutyCycle() {
	inst.cancelChannel <- true
	if inst.onTicker != nil {
		inst.onTicker.Stop()
	}
	if inst.offTicker != nil {
		inst.offTicker.Stop()
	}
	if inst.offDelayTimer != nil {
		inst.offDelayTimer.Stop()
	}
	inst.lastEnable = false
	inst.WritePinFalse(node.Out)
}
