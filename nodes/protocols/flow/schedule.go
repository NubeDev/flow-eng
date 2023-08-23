package flow

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
	log "github.com/sirupsen/logrus"
	"time"
)

type FFSchedule struct {
	*node.Spec
	topic       string
	lastPayload *covPayload
	lastValue   bool
	lastUpdate  time.Time
	hasWritten  bool
}

func NewFFSchedule(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowSchedule, category)
	inputs := node.BuildInputs()
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	lastUpdated := node.BuildOutput(node.LastUpdated, node.TypeString, nil, body.Outputs)
	payload := node.BuildOutput(node.OutPayload, node.TypeFloat, nil, body.Outputs)
	periodStart := node.BuildOutput(node.PeriodStart, node.TypeString, nil, body.Outputs)
	periodStop := node.BuildOutput(node.PeriodStop, node.TypeString, nil, body.Outputs)
	periodStartUnix := node.BuildOutput(node.PeriodStartUnix, node.TypeFloat, nil, body.Outputs)
	periodStopUnix := node.BuildOutput(node.PeriodStopUnix, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out, lastUpdated, payload, periodStart, periodStop, periodStartUnix, periodStopUnix)
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body = node.SetNoParent(body)
	pnt := &FFSchedule{body, "", nil, false, time.Now(), false}
	return pnt, nil
}

func (inst *FFSchedule) getSchedules() ([]*Schedule, []string, error) {
	s := inst.GetStore()
	var scheduleNames []string
	if s == nil {
		return nil, scheduleNames, errors.New("failed to find store")
	}
	parentId := inst.GetParentId()
	topic := fmt.Sprintf("schedulesList_%s", parentId)
	d, ok := s.Get(topic)
	var data []*Schedule
	if ok {
		data = d.([]*Schedule)
		for _, datum := range data {
			scheduleNames = append(scheduleNames, datum.Name)
		}
	}
	return data, scheduleNames, nil
}

func (inst *FFSchedule) getResult() {
	settings, err := getScheduleSettings(inst.GetSettings())
	if err != nil {
		log.Errorf("Flow Network Schedules getResult() err: %s", err.Error())
	}
	schedules, _, err := inst.getSchedules()
	if err != nil {
		log.Errorf("Flow Network Schedules getResult() err: %s", err.Error())
	}
	inst.hasWritten = false
	for _, schedule := range schedules {
		if settings.Schedule == schedule.Name {
			value := schedule.IsActive
			if inst.lastValue != value {
				inst.lastValue = value
				inst.lastUpdate = time.Now()
			}
			inst.hasWritten = true
			inst.WritePinBool(node.Out, value)
			inst.WritePin(node.LastUpdated, ttime.TimeSince(inst.lastUpdate))
			inst.WritePinFloat(node.OutPayload, schedule.Payload)
			inst.WritePin(node.PeriodStart, schedule.PeriodStartString)
			inst.WritePin(node.PeriodStop, schedule.PeriodStopString)
			inst.WritePinFloat(node.PeriodStartUnix, float64(schedule.PeriodStart))
			inst.WritePinFloat(node.PeriodStopUnix, float64(schedule.PeriodStop))
			inst.SetSubTitle(schedule.Name)
			inst.SetWaringIcon(string(emoji.GreenCircle))
			inst.hasWritten = true
		}
	}
}

func (inst *FFSchedule) Process() {
	loopCount, _ := inst.Loop()
	if loopCount == 5 {
		inst.getResult()
	} else if loopCount%50 == 0 {
		inst.getResult()
	}
	if !inst.hasWritten {
		inst.WritePinNull(node.Out)
		inst.WritePinNull(node.LastUpdated)
		inst.WritePinNull(node.OutPayload)
		inst.WritePinNull(node.NextStart)
		inst.WritePinNull(node.NextStop)
		inst.SetWaringIcon(string(emoji.OrangeCircle))
	}

}

func (inst *FFSchedule) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
