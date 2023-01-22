package timing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/str"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"time"
)

type DutyCycle struct {
	*node.Spec
	onTicker      *time.Ticker
	offTimer      *time.Timer
	cancelChannel chan bool
	lastInterval  time.Duration
	lastDuty      float64
	lastEnable    bool
}

func NewDutyCycle(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, dutyCycle, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, nil)
	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, str.New("interval"))
	dutyCycleInput := node.BuildInput(node.DutyCycle, node.TypeFloat, nil, body.Inputs, str.New("duty_cycle"))
	inputs := node.BuildInputs(enable, interval, dutyCycleInput)
	out := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	node := &DutyCycle{body, nil, nil, nil, 10, 50, false}
	node.SetSchema(node.buildSchema())
	return node, nil
}

func (inst *DutyCycle) Process() {

	enable, _ := inst.ReadPinAsBool(node.Enable)
	if !enable {
		inst.disableDutyCycle()
		inst.lastEnable = false
		inst.WritePinFalse(node.Outp)
		return
	}

	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	if intervalDuration <= 0 {
		intervalDuration = 10 * time.Second
	}
	if intervalDuration != inst.lastInterval {
		inst.setSubtitle(intervalDuration)
	}

	dutyCyclePerc := inst.ReadPinOrSettingsFloat(node.DutyCycle)
	if dutyCyclePerc < 0 || dutyCyclePerc > 100 {
		dutyCyclePerc = 50
	}

	// Check if there are settings that require a restart
	if enable && !inst.lastEnable || intervalDuration != inst.lastInterval || dutyCyclePerc != inst.lastDuty {
		inst.restartDutyCycle(intervalDuration, dutyCyclePerc)
	}
	inst.lastInterval = intervalDuration
	inst.lastDuty = dutyCyclePerc
	inst.lastEnable = enable

}

func (inst *DutyCycle) restartDutyCycle(intervalDuration time.Duration, dutyCycle float64) error {
	inst.disableDutyCycle() // stop existing timers

	if intervalDuration <= 0 || (dutyCycle < 0 || dutyCycle > 100) {
		return errors.New("restartDutyCycle() err: invalid inputs")
	}

	delayBetweenOnAndOff := intervalDuration.Seconds() * (dutyCycle / 100)
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
	inst.WritePinTrue(node.Outp)
	inst.offTimer = time.AfterFunc(delayBetweenOnAndOffDuration, func() {
		inst.WritePinFalse(node.Outp)
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

func (inst *DutyCycle) setSubtitle(intervalDuration time.Duration) {
	subtitleText := intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type DutyCycleSettingsSchema struct {
	Interval  schemas.Number     `json:"interval"`
	TimeUnits schemas.EnumString `json:"interval_time_units"`
	DutyCycle schemas.Number     `json:"duty_cycle"`
}

type DutyCycleSettings struct {
	Interval  float64 `json:"interval"`
	TimeUnits string  `json:"interval_time_units"`
	DutyCycle float64 `json:"duty_cycle"`
}

func (inst *DutyCycle) buildSchema() *schemas.Schema {
	props := &DutyCycleSettingsSchema{}
	// time selection
	props.Interval.Title = "Period"
	props.Interval.Default = 1

	// time selection
	props.TimeUnits.Title = "Period Units"
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
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "interval_time_units", "duty_cycle"},
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
