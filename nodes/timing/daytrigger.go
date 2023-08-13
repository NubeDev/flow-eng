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
	log "github.com/sirupsen/logrus"
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
	timeSetting := node.BuildInput(node.Time, node.TypeString, 0, body.Inputs, true, false)
	inputs := node.BuildInputs(enable, timeSetting)

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
	inst.at = inst.ReadPinOrSettingsString(node.Time)
	day := settings.DaySelection

	_, _, _, err := ttime.ParseTime(inst.at)
	if err != nil {
		log.Error(fmt.Sprintf("day-trigger: failed to parse time err: %s", err.Error()))
		return
	}

	scheduler := gocron.NewScheduler()

	if day == everyDay {
		scheduler.Every(1).Day().At(inst.at).Do(inst.job)
	}
	if day == ttime.Sun {
		scheduler.Every(1).Sunday().At(inst.at).Do(inst.job)
	}
	if day == ttime.Mon {
		scheduler.Every(1).Monday().At(inst.at).Do(inst.job)
	}
	if day == ttime.Tue {
		scheduler.Every(1).Tuesday().At(inst.at).Do(inst.job)
	}
	if day == ttime.Wed {
		scheduler.Every(1).Wednesday().At(inst.at).Do(inst.job)
	}
	if day == ttime.Thur {
		scheduler.Every(1).Thursday().At(inst.at).Do(inst.job)
	}
	if day == ttime.Fri {
		scheduler.Every(1).Friday().At(inst.at).Do(inst.job)
	}
	if day == ttime.Sat {
		scheduler.Every(1).Saturday().At(inst.at).Do(inst.job)
	}
	scheduler.Start()
	inst.setSubtitle(day)
}

func (inst *DayTrigger) setSubtitle(day string) {
	var title string
	fmt.Println(inst.at)
	if day == everyDay {
		title = fmt.Sprintf("trigger every day at:(%s) hold for: (%d:sec)", inst.at, inst.lockDuration)
	} else {
		title = fmt.Sprintf("trigger at:(%s:%s) hold for: (%d:sec)", inst.at, day, inst.lockDuration)
	}

	inst.SetSubTitle(title)
}

type dayTriggerSettingsSchema struct {
	Time         schemas.String     `json:"time"`
	DaySelection schemas.EnumString `json:"day_selection"`
	LockDuration schemas.Integer    `json:"lock_duration"`
}

type dayTriggerSettings struct {
	Time         string `json:"time"`
	DaySelection string `json:"day_selection"`
	LockDuration int    `json:"lock_duration"`
}

const everyDay = "everyDay"

func (inst *DayTrigger) buildSchema() *schemas.Schema {
	props := &dayTriggerSettingsSchema{}

	// time selection
	props.DaySelection.Title = "Day To Trigger"
	props.DaySelection.Default = everyDay
	props.DaySelection.Options = []string{everyDay, ttime.Sun, ttime.Mon, ttime.Tue, ttime.Wed, ttime.Thur, ttime.Fri, ttime.Sat}
	props.DaySelection.EnumName = []string{everyDay, ttime.Sun, ttime.Mon, ttime.Tue, ttime.Wed, ttime.Thur, ttime.Fri, ttime.Sat}

	props.Time.Title = "trigger at (eg 08:00, 20:00)"
	props.Time.Default = "08:00"

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
		"ui:order": array.Slice{"time", "day_selection", "lock_duration"},
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
