package timing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
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
	node := &DutyCycle{body, nil, nil, nil, 10, 50, false}
	node.SetSchema(node.buildSchema())
	return node, nil
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
	inst.offTimer = time.AfterFunc(delayBetweenOnAndOffDuration, func() {
		inst.WritePinFalse(node.Out)
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

// Custom Node Settings Schema

type DutyCycleSettingsSchema struct {
	Interval  schemas.Number     `json:"interval"`
	TimeUnits schemas.EnumString `json:"time_units"`
	DutyCycle schemas.Number     `json:"duty_cycle"`
}

type DutyCycleSettings struct {
	Interval  float64 `json:"interval"`
	TimeUnits string  `json:"time_units"`
	DutyCycle float64 `json:"duty_cycle"`
}

func (inst *DutyCycle) buildSchema() *schemas.Schema {
	props := &DutyCycleSettingsSchema{}
	// time selection
	props.Interval.Title = "Interval"
	props.Interval.Default = 1

	// time selection
	props.TimeUnits.Title = "Time Units"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.TimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// retrigger selection
	props.DutyCycle.Title = "Duty Cycle Percent"
	props.DutyCycle.Default = 50
	props.DutyCycle.Minimum = 0
	props.DutyCycle.Maximum = 100

	schema.Set(props)

	uiSchema := array.Map{
		"time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "time_units", "duty_cycle"},
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

func (inst *DutyCycle) getSettings(body map[string]interface{}) (*DutyCycleSettings, error) {
	settings := &DutyCycleSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
