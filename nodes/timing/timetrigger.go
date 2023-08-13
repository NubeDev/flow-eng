package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/carlescere/scheduler"
	"time"
)

type TimeTrigger struct {
	*node.Spec
	lock         bool
	every        int
	lockDuration int
}

func NewTimeTrigger(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, timeTrigger, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, true)

	inputs := node.BuildInputs(enable)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &TimeTrigger{body, false, 0, 0}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *TimeTrigger) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.init()
	}
	inst.WritePinBool(node.Out, inst.lock)

}

func (inst *TimeTrigger) job() {
	inst.lock = true
	if inst.every <= 1 {
		time.Sleep(800 * time.Millisecond)
	} else {
		time.Sleep(time.Duration(inst.lockDuration) * time.Second)
	}
	inst.lock = false
}

func (inst *TimeTrigger) init() {
	settings, _ := inst.getSettings(inst.GetSettings())
	inst.every = settings.Trigger
	inst.lockDuration = settings.LockDuration
	timeUnits := settings.TimeUnits
	if timeUnits == ttime.Sec {
		scheduler.Every(inst.every).Seconds().Run(inst.job)
	}
	if timeUnits == ttime.Min {
		scheduler.Every(inst.every).Minutes().Run(inst.job)
	}
	if timeUnits == ttime.Hr {
		scheduler.Every(inst.every).Hours().Run(inst.job)
	}
	if timeUnits == ttime.Day {
		scheduler.Every(inst.every).Day().Run(inst.job)
	}
	inst.setSubtitle(timeUnits)
}

func (inst *TimeTrigger) setSubtitle(timeUnits string) {
	title := fmt.Sprintf("trigger at:(%d:%s) lock: (%d:%s)", inst.every, timeUnits, inst.lockDuration, timeUnits)
	inst.SetSubTitle(title)
}

type schedulerSettingsSchema struct {
	Trigger      schemas.Integer    `json:"trigger"`
	LockDuration schemas.Integer    `json:"lock_duration"`
	TimeUnits    schemas.EnumString `json:"time_units"`
}

type schedulerSettings struct {
	Trigger      int    `json:"trigger"`
	LockDuration int    `json:"lock_duration"`
	TimeUnits    string `json:"time_units"`
}

func (inst *TimeTrigger) buildSchema() *schemas.Schema {
	props := &schedulerSettingsSchema{}

	props.Trigger.Title = "trigger at"
	props.Trigger.Default = 2
	props.Trigger.Minimum = 2
	props.Trigger.Maximum = 100000000

	// time selection
	props.TimeUnits.Title = "Interval Units"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Sec, ttime.Min, ttime.Hr, ttime.Day}
	props.TimeUnits.EnumName = []string{ttime.Sec, ttime.Min, ttime.Hr, ttime.Day}

	props.LockDuration.Title = "for duration (x) seconds"
	props.LockDuration.Default = 1
	props.LockDuration.Minimum = 1
	props.LockDuration.Maximum = 100000000

	schema.Set(props)

	uiSchema := array.Map{
		"time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"trigger", "time_units", "lock_duration"},
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
func (inst *TimeTrigger) getSettings(body map[string]interface{}) (*schedulerSettings, error) {
	settings := &schedulerSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
