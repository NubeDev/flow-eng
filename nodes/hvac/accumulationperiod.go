package hvac

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/gorhill/cronexpr"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
)

type AccumulationPeriod struct {
	*node.Spec
	lastAccumulation          float64
	lastPeriodEndAccumulation float64
	lastPeriodConsumption     float64
	cron                      *cron.Cron
	cronExp                   string
}

func NewAccumulationPeriod(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, accumulationPeriod, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs, false)
	input := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false)
	inputs := node.BuildInputs(enable, input)

	periodConsumption := node.BuildOutput(node.PeriodConsumption, node.TypeFloat, nil, body.Outputs)
	lastAccum := node.BuildOutput(node.LastPeriodEndVal, node.TypeFloat, nil, body.Outputs)
	periodDuration := node.BuildOutput(node.PeriodDuration, node.TypeFloat, nil, body.Outputs)
	nextTrigger := node.BuildOutput(node.NextTrigger, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(periodConsumption, lastAccum, periodDuration, nextTrigger)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &AccumulationPeriod{body, 0, 0, 0, nil, ""}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *AccumulationPeriod) Process() {
	if inst.cron == nil {
		settings, _ := inst.getSettings(inst.GetSettings())
		period := settings.AccumulationPeriod
		cronExp := ""
		if period <= 60 {
			cronExp = "*/" + strconv.Itoa(period) + " * * * *"
		} else {
			cronExp = "0 */" + strconv.Itoa(period/60) + " * * *"
		}
		inst.cronExp = cronExp
		inst.cron = cron.New()
		inst.cron.AddFunc(cronExp, inst.calculateAccumulation)
		inst.cron.Start()
		inst.WritePinFloat(node.PeriodDuration, float64(period))
		nextTrigger := cronexpr.MustParse(inst.cronExp).Next(time.Now())
		inst.WritePin(node.NextTrigger, nextTrigger)
		input, inNull := inst.ReadPinAsFloat(node.In)
		if !inNull {
			inst.lastAccumulation = input
		}
		inst.WritePinFloat(node.PeriodConsumption, inst.lastPeriodConsumption)
		inst.WritePinFloat(node.LastPeriodEndVal, inst.lastPeriodEndAccumulation)
	}
}

func (inst *AccumulationPeriod) calculateAccumulation() {
	input, inNull := inst.ReadPinAsFloat(node.In)
	if !inNull {
		inst.lastAccumulation = input
	}

	inst.lastPeriodConsumption = inst.lastAccumulation - inst.lastPeriodEndAccumulation
	inst.lastPeriodEndAccumulation = inst.lastAccumulation
	inst.WritePinFloat(node.PeriodConsumption, inst.lastPeriodConsumption)
	inst.WritePinFloat(node.LastPeriodEndVal, inst.lastPeriodEndAccumulation)
	nextTrigger := cronexpr.MustParse(inst.cronExp).Next(time.Now())
	inst.WritePin(node.NextTrigger, nextTrigger)
}

func (inst *AccumulationPeriod) Stop() {
	inst.cron.Stop()
}

// Custom Node Settings Schema

type AccumulationPeriodSettingsSchema struct {
	AccumulationPeriod schemas.EnumInt `json:"accumulation-period"`
}

type AccumulationPeriodSettings struct {
	AccumulationPeriod int `json:"accumulation-period"`
}

func (inst *AccumulationPeriod) buildSchema() *schemas.Schema {
	props := &AccumulationPeriodSettingsSchema{}

	// accumulation period
	props.AccumulationPeriod.Title = "Accumulation Period (minutes)"
	props.AccumulationPeriod.Default = 15
	props.AccumulationPeriod.Options = []int{1, 2, 3, 4, 5, 6, 10, 12, 15, 20, 30, 60, 120, 180, 240, 360, 720}
	props.AccumulationPeriod.EnumName = []int{1, 2, 3, 4, 5, 6, 10, 12, 15, 20, 30, 60, 120, 180, 240, 360, 720}

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"accumulation-period"},
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

func (inst *AccumulationPeriod) getSettings(body map[string]interface{}) (*AccumulationPeriodSettings, error) {
	settings := &AccumulationPeriodSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
