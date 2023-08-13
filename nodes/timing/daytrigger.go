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

type DayTrigger struct {
	*node.Spec
	lock         bool
	day          string
	at           string
	lockDuration int
}

func NewDayTrigger(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, dayTrigger, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false, true)

	inputs := node.BuildInputs(enable)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &DayTrigger{body, false, "", "", 0}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *DayTrigger) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		inst.init()
	}
	inst.WritePinBool(node.Out, inst.lock)

}

func (inst *DayTrigger) job() {
	inst.lock = true
	time.Sleep(time.Duration(inst.lockDuration) * time.Second)
	inst.lock = false
}

func (inst *DayTrigger) init() {
	settings, _ := inst.getSettings(inst.GetSettings())
	inst.lockDuration = settings.LockDuration
	day := settings.DaySelection
	if day == ttime.Sun {
		scheduler.Every().Sunday().At(inst.at).Run(inst.job)
	}
	if day == ttime.Mon {
		scheduler.Every().Monday().At(inst.at).Run(inst.job)
	}
	if day == ttime.Tue {
		scheduler.Every().Tuesday().At(inst.at).Run(inst.job)
	}
	if day == ttime.Wed {
		scheduler.Every().Wednesday().At(inst.at).Run(inst.job)
	}
	if day == ttime.Thur {
		scheduler.Every().Thursday().At(inst.at).Run(inst.job)
	}
	if day == ttime.Friday {
		scheduler.Every().Friday().At(inst.at).Run(inst.job)
	}
	if day == ttime.Saturday {
		scheduler.Every().Saturday().At(inst.at).Run(inst.job)
	}
	inst.setSubtitle(day)
}

func (inst *DayTrigger) setSubtitle(timeUnits string) {
	title := fmt.Sprintf("trigger at:(%d:%s) lock: (%d:%s)", inst.at, timeUnits, inst.lockDuration, timeUnits)
	inst.SetSubTitle(title)
}

type dayTriggerSettingsSchema struct {
	TriggerAt    schemas.String     `json:"trigger"`
	DaySelection schemas.EnumString `json:"day_selection"`
	LockDuration schemas.Integer    `json:"lock_duration"`
}

type dayTriggerSettings struct {
	TriggerAt    string `json:"trigger"`
	DaySelection string `json:"day_selection"`
	LockDuration int    `json:"lock_duration"`
}

func (inst *DayTrigger) buildSchema() *schemas.Schema {
	props := &dayTriggerSettingsSchema{}

	// time selection
	props.DaySelection.Title = "Day To Trigger"
	props.DaySelection.Default = ttime.Mon
	props.DaySelection.Options = []string{ttime.Sun, ttime.Mon, ttime.Tue, ttime.Wed, ttime.Thur, ttime.Fri, ttime.Sat}
	props.DaySelection.EnumName = []string{ttime.Sun, ttime.Mon, ttime.Tue, ttime.Wed, ttime.Thur, ttime.Fri, ttime.Sat}

	props.TriggerAt.Title = "trigger at (eg 08:00, 20:00)"
	props.TriggerAt.Default = "08:00"

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
		"ui:order": array.Slice{"trigger", "day_selection", "lock_duration"},
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
func (inst *DayTrigger) getSettings(body map[string]interface{}) (*dayTriggerSettings, error) {
	settings := &dayTriggerSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
