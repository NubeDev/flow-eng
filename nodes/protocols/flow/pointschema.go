package flow

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"strings"
)

func pointTopic(selected string) string {
	parts := strings.Split(selected, ":")
	if len(parts) >= 3 {
		return fmt.Sprintf("+/+/+/+/+/+/rubix/points/value/cov/all/%s/+/%s/+/%s/+/%s", parts[0], parts[1], parts[2], parts[3])
		// na/na/na/na/Ugbo2JjAnfc939rRBTpuoT/na/rubix/points/value/cov/all/system/net_ebc64b5754674378/net/dev_21223294a79742a1/dev/pnt_94ea3ea254dc440a/pnt
	}
	return ""
}

func getPoints(points []*point) (names []string) {
	for _, p := range points {
		names = append(names, p.Name)
	}
	return names
}

type pointNodeSchema struct {
	Point schemas.EnumString `json:"point"`
}

func (inst *Point) buildSchema() *schemas.Schema {
	s := inst.GetStore()
	if s == nil {
		return nil
	}
	data, ok := s.Get(fmt.Sprintf("pointsList_%s", inst.GetParentId()))
	if !ok {
		return nil
	}
	d, _ := data.([]*point)
	names := getPoints(d)
	props := &pointNodeSchema{}
	props.Point.Title = "point"
	if len(names) > 0 {
		props.Point.Default = names[0]
	} else {
		names = append(names, "no connection has been added")
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
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

func (inst *Point) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}