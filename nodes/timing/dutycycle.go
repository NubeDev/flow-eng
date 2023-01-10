package timing

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
	"time"
)

type DutyCycle struct {
	*node.Spec
	onTicker         *time.Ticker
	offTimer         *time.Timer
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
	body.SetSchema(buildDefaultSchema())
	return &DutyCycle{body, nil, nil, nil, 10, 50, false}, nil
}

func (inst *DutyCycle) Process() {
	// settings, _ := getSettings(inst.GetSettings())

	enable, _ := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.disableDutyCycle()
		inst.lastEnable = false
		inst.WritePinFalse(node.Out)
		return
	}

	intervalSecs, _ := inst.ReadPinAsFloat(node.IntervalSecs)
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
	inst.disableDutyCycle() // stop existing timers

	if intervalSeconds <= 0 || (dutyCycle < 0 || dutyCycle > 100) {
		return errors.New("restartDutyCycle() err: invalid inputs")
	}
	intervalDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", intervalSeconds))

	delayBetweenOnAndOff := intervalSeconds * (dutyCycle / 100)
	delayBetweenOnAndOffDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", delayBetweenOnAndOff))

	cancel := make(chan bool)
	inst.cancelChannel = cancel
	inst.onTicker = time.NewTicker(intervalDuration)
	inst.startIteration(delayBetweenOnAndOffDuration)

	go func() {
		for {
			select {
			case <-cancel:
				return
			case <-inst.onTicker.C:
				inst.startIteration(delayBetweenOnAndOffDuration)
			}
		}
	}()

	return nil
}

func (inst *DutyCycle) startIteration(delayBetweenOnAndOffDuration time.Duration) {
	inst.WritePinTrue(node.Out)
	log.Infof("DutyCycle: ON")
	inst.offTimer = time.AfterFunc(delayBetweenOnAndOffDuration, func() {
		inst.WritePinFalse(node.Out)
		log.Infof("DutyCycle: OFF")
	})
}

func (inst *DutyCycle) disableDutyCycle() {
	// log.Infof("disableDutyCycle()")
	if inst.cancelChannel != nil {
		inst.cancelChannel <- true
		inst.cancelChannel = nil
	}
	if inst.onTicker != nil {
		inst.onTicker.Stop()
		inst.onTicker = nil
	}
	if inst.offTimer != nil {
		inst.offTimer.Stop()
		inst.offTimer = nil
	}
}

func (inst *DutyCycle) Stop() {
	inst.disableDutyCycle()
}
