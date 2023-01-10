package timing

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"strings"
	"time"
)

type OneShot struct {
	*node.Spec
	timer        *time.Timer
	outputActive bool
	lastIn       bool
	lastReset    bool
}

type OneShotSchema struct {
	Time      schemas.EnumString `json:"time"`
	Duration  schemas.Number     `json:"duration"`
	Retrigger schemas.Boolean    `json:"retrigger"`
}

type OneShotSettings struct {
	Time      string        `json:"time"`
	Duration  time.Duration `json:"duration"`
	Retrigger bool          `json:"retrigger"`
}

func (inst *OneShot) getSettings(body map[string]interface{}) (*nodeSettings, error) {
	settings := &nodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

func (inst *OneShot) buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	// time selection
	props.Duration.Title = "duration"
	props.Duration.Default = 1

	// time selection
	props.Time.Title = "time"
	props.Time.Default = ttime.Sec
	props.Time.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.Time.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	pprint.PrintJSON(props)
	schema.Set(props)

	fmt.Println(fmt.Sprintf("buildSchema() props: %+v", props))
	pprint.PrintJSON(props)

	uiSchema := array.Map{
		"time": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Set delay time",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	fmt.Println(fmt.Sprintf("buildSchema() s: %+v", s))
	pprint.PrintJSON(s)
	return s
}

func NewOneShot(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, oneShot, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)       // TODO: this input shouldn't have a manual override value
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(in, reset)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &OneShot{body, nil, false, true, true}, nil
}

func (inst *OneShot) Process() {
	retrigger := false

	settings, _ := getSettings(inst.GetSettings())
	if settings != nil {
		t := strings.Replace(settings.Duration.String(), "ns", "", -1)
		inst.SetSubTitle(fmt.Sprintf("setting: %s %s", t, settings.Time))
	}

	in, _ := inst.ReadPinAsBool(node.In)
	if in && !inst.lastIn {
		if retrigger || !inst.outputActive {
			oneShotDuration := ttime.Duration(settings.Duration, settings.Time)
			inst.StartOneShot(oneShotDuration)
		}
	}
	inst.lastIn = in

	reset, _ := inst.ReadPinAsBool(node.Reset)
	if reset && !inst.lastReset {
		if inst.outputActive {
			inst.StopOneShotTimer(true)
		}
	}
	inst.lastReset = reset
}

func (inst *OneShot) StartOneShot(duration time.Duration) {
	if inst.timer != nil {
		inst.StopOneShotTimer(false)
	}
	inst.timer = time.AfterFunc(duration, func() {
		inst.WritePinFalse(node.Out)
		inst.outputActive = false
		inst.timer = nil
	})
	inst.WritePinTrue(node.Out)
	inst.outputActive = true
}

func (inst *OneShot) StopOneShotTimer(reset bool) {
	if inst.timer != nil {
		inst.timer.Stop()
	}
	if reset {
		inst.WritePinFalse(node.Out)
		inst.outputActive = false
	}
}

func (inst *OneShot) Start() {
	inst.WritePinFalse(node.Out)
	inst.outputActive = false
}

func (inst *OneShot) Stop() {
	inst.StopOneShotTimer(true)
}
