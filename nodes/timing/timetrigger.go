package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/jasonlvhit/gocron"
	"time"
)

const timeTriggerDesc = `node used for triggering an true/false value as per the node settings`

type TimedTrigger struct {
	*node.Spec
	lock         bool
	every        uint64
	lockDuration int
}

func NewTimedTrigger(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, timedTrigger, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, true, body.Inputs, false, false, node.SetInputHelp(node.EnableHelp))
	inputs := node.BuildInputs(enable)
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs, node.SetOutputHelp(node.OutHelp))
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	n := &TimedTrigger{body, false, 0, 0}
	n.SetSchema(n.buildSchema())
	n.SetHelp(timeTriggerDesc)
	return n, nil
}

func (inst *TimedTrigger) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.init()
	}
	inst.WritePinBool(node.Out, inst.lock)

}

func (inst *TimedTrigger) job() {
	inst.lock = true
	if inst.every <= 1 {
		time.Sleep(800 * time.Millisecond)
	} else {
		time.Sleep(time.Duration(inst.lockDuration) * time.Second)
	}
	inst.lock = false
}

func (inst *TimedTrigger) init() {
	settings, _ := inst.getSettings(inst.GetSettings())
	inst.every = settings.Trigger
	inst.lockDuration = settings.LockDuration
	timeUnits := settings.TimeUnits
	scheduler := gocron.NewScheduler()
	if timeUnits == ttime.Sec {
		scheduler.Every(inst.every).Seconds().Do(inst.job)
	}
	if timeUnits == ttime.Min {
		scheduler.Every(inst.every).Minutes().Do(inst.job)
	}
	if timeUnits == ttime.Hr {
		scheduler.Every(inst.every).Hours().Do(inst.job)
	}
	if timeUnits == ttime.Day {
		scheduler.Every(inst.every).Day().Do(inst.job)
	}
	inst.setSubtitle(timeUnits)
	scheduler.Start()
}

func (inst *TimedTrigger) setSubtitle(timeUnits string) {
	title := fmt.Sprintf("trigger at:(%d:%s) lock: (%d:%s)", inst.every, timeUnits, inst.lockDuration, timeUnits)
	inst.SetSubTitle(title)
}

type schedulerSettingsSchema struct {
	Trigger      schemas.Integer    `json:"trigger"`
	LockDuration schemas.Integer    `json:"lock_duration"`
	TimeUnits    schemas.EnumString `json:"time_units"`
}

type schedulerSettings struct {
	Trigger      uint64 `json:"trigger"`
	LockDuration int    `json:"lock_duration"`
	TimeUnits    string `json:"time_units"`
}

func (inst *TimedTrigger) buildSchema() *schemas.Schema {
	props := &schedulerSettingsSchema{}

	props.Trigger.Title = "trigger at"
	props.Trigger.Help = "this is help"
	props.Trigger.Default = 2
	props.Trigger.Minimum = 2
	props.Trigger.Maximum = 100000000

	// time selection
	props.TimeUnits.Title = "Interval Units"
	props.TimeUnits.Help = "this is help"
	props.TimeUnits.Default = ttime.Sec
	props.TimeUnits.Options = []string{ttime.Sec, ttime.Min, ttime.Hr, ttime.Day}
	props.TimeUnits.EnumName = []string{ttime.Sec, ttime.Min, ttime.Hr, ttime.Day}

	props.LockDuration.Title = "for duration (x) seconds"
	props.LockDuration.Help = "this is help"
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
func (inst *TimedTrigger) getSettings(body map[string]interface{}) (*schedulerSettings, error) {
	settings := &schedulerSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
