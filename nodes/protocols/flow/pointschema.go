package flow

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"strings"
)

func fixTopic(topic string) (fixedTopic, pointUUID string) {
	parts := strings.Split(topic, "/")
	if len(parts) == 12 {
		pointUUID = parts[10]
		parts[6] = "+"
		parts[8] = "+"
		parts[10] = "+"
		return strings.Join(parts, "/"), pointUUID
	}
	return "", ""
}

func makePointTopic(selected string) string {
	parts := strings.Split(selected, ":")
	if len(parts) >= 3 {
		return fmt.Sprintf("rubix/points/value/cov/all/%s/+/%s/+/%s/+/%s", parts[0], parts[1], parts[2], parts[3])
	}
	return ""
}

const selectAPoint = "Select a point"

func getPoints(points []*point) (names []string) {
	names = append(names, selectAPoint)
	for _, p := range points {
		names = append(names, p.Name)
	}
	return names
}

type pointNodeSchema struct {
	Point schemas.EnumString `json:"point"`
}

func (inst *FFPoint) buildSchema() *schemas.Schema {
	s := inst.GetStore()
	var data interface{}
	var ok bool
	var names []string
	if s != nil {
		data, ok = s.Get(fmt.Sprintf("pointsList_%s", inst.GetParentId()))
		if ok {
			d, _ := data.([]*point)
			names = getPoints(d)
		}
	} else {
		names = getPoints(nil)
	}
	props := &pointNodeSchema{}
	props.Point.Title = "point"
	if len(names) > 0 {
		props.Point.Default = names[0]
	}
	props.Point.Options = names
	props.Point.EnumName = names
	schema.Set(props)
	sch := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return sch
}

func (inst *FFPointWrite) buildSchema() *schemas.Schema {
	s := inst.GetStore()
	var data interface{}
	var ok bool
	var names []string
	if s != nil {
		data, ok = s.Get(fmt.Sprintf("pointsList_%s", inst.GetParentId()))
		if ok {
			d, _ := data.([]*point)
			names = getPoints(d)
		}
	} else {
		names = getPoints(nil)
	}
	props := &pointNodeSchema{}
	props.Point.Title = "point"
	props.Point.Title = "point"
	if len(names) > 0 {
		props.Point.Default = names[0]
	}
	props.Point.Options = names
	props.Point.EnumName = names
	schema.Set(props)
	sch := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return sch
}

type pointSettings struct {
	Point string `json:"point"`
}

func getPointSettings(body map[string]interface{}) (*pointSettings, error) {
	settings := &pointSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &settings)
	if err != nil {
		return nil, err
	}
	return settings, err
}

func getScheduleSettings(body map[string]interface{}) (*scheduleSettings, error) {
	settings := &scheduleSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &settings)
	if err != nil {
		return nil, err
	}
	return settings, err
}

type scheduleNodeSchema struct {
	Name schemas.EnumString `json:"schedule"`
}

type scheduleSettings struct {
	Name string `json:"schedule"`
}

func (inst *FFSchedule) buildSchema() *schemas.Schema {
	_, names, _ := inst.getSchedules()
	props := &scheduleNodeSchema{}
	props.Name.Title = "select a schedules"
	props.Name.Options = names
	props.Name.EnumName = names
	schema.Set(props)
	sch := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return sch
}
