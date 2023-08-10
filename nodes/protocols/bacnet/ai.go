package bacnetio

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
)

type AI struct {
	*node.Spec
	objectID      points.ObjectID
	objectType    points.ObjectType
	pointUUID     string
	store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

func NewAI(body *node.Spec, opts ...any) (node.Node, error) {
	bn := bacnetOpts(opts...)
	body = node.Defaults(body, bacnetAI, Category)
	var inputs []*node.Input // no inputs required
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	flowOptions := &toFlowOptions{}
	n := &AI{
		body,
		0,
		points.AnalogInput,
		"",
		bn.Store,
		bn.Application,
		bn.MqttClient,
		flowOptions,
	}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *AI) Process() {
	settings, err := inst.getSettings(inst.GetSettings())
	_, firstLoop := inst.Loop()
	if firstLoop {
		transformProps := inst.getTransformProps(settings)
		inst.setObjectId(settings)
		ioType := settings.IoType
		if ioType == "" {
			ioType = string(points.IoTypeVolts)
		}
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		point := addPoint(points.IoType(ioType), objectType, inst.objectID, isWriteable, isIO, true, inst.application, transformProps)
		name := inst.GetNodeName()
		parentTopic := helpers.CleanParentName(name, inst.GetParentName())
		if parentTopic != "" {
			name = parentTopic
		}
		point.Name = name
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d err:%s", objectType, inst.objectID, err.Error())
			inst.WritePinNull(node.Out)
			return
		}
		s := inst.GetStore()
		if s == nil {
			log.Errorf("bacnet-server add new point failed to get store type:%s-%d err:%s", objectType, inst.objectID, err.Error())
			inst.WritePinNull(node.Out)
			return
		}
		point, err = inst.store.AddPoint(point, true)
		if err != nil {
			log.Errorf("bacnet-server add new point type: %s-%d", objectType, inst.objectID)
		} else {
			s.Set(setUUID(inst.GetParentId(), points.AnalogInput, inst.objectID), point, 0)
		}
	}

	pv, modbusUpdated, err := inst.getPV(points.AnalogInput, inst.objectID)
	if err != nil {
		inst.WritePinNull(node.Out)
		return
	}
	if pv == nil {
		inst.WritePinNull(node.Out)
	} else if modbusUpdated {
		inst.WritePinFloat(node.Out, *pv, settings.Decimal)
	} else {
		inst.WritePinNull(node.Out)
	}
}

func (inst *AI) setObjectId(settings *AISettings) {
	devNum := settings.DeviceNumber
	uiNum := settings.InputNumber
	id := ((devNum - 1) * 8) + uiNum
	inst.objectID = points.ObjectID(id)
	name := bacnetAddress(4, "AI", "UI")
	if len(name) >= id {
		if settings != nil {
			ioType := strings.ReplaceAll(settings.IoType, "_", " ")
			inst.SetSubTitle(strings.ToUpper(fmt.Sprintf("%s %s", name[id-1], ioType)))
		} else {
			inst.SetSubTitle(name[id-1])
		}
	}
}

func (inst *AI) getPV(objType points.ObjectType, id points.ObjectID) (*float64, bool, error) {
	pnt, ok := inst.getPoint(objType, id)
	// fmt.Println(fmt.Sprintf("AI getPV() pnt.PresentValue: %v", pnt.PresentValue))
	if ok {
		return pnt.PresentValue, pnt.ModbusUpdated, nil
	}
	return float.New(0), false, nil
}

func (inst *AI) getPoint(objType points.ObjectType, id points.ObjectID) (*points.Point, bool) {
	s := inst.GetStore()
	if s == nil {
		return nil, false
	}
	d, ok := s.Get(setUUID(inst.ParentId, objType, id))
	if ok {
		return d.(*points.Point), true
	}
	return nil, false
}

// Custom Node Settings Schema

type AISettingsSchema struct {
	DeviceNumber schemas.Integer        `json:"device-number"`
	InputNumber  schemas.Integer        `json:"input-number"`
	IoType       schemas.EnumString     `json:"io-type"`
	Decimal      schemas.Number         `json:"decimal"`
	ScaleEnable  schemas.Boolean        `json:"scale-enable"`
	ScaleInMin   schemas.NumberNoLimits `json:"scale-in-min"`
	ScaleInMax   schemas.NumberNoLimits `json:"scale-in-max"`
	ScaleOutMin  schemas.NumberNoLimits `json:"scale-out-min"`
	ScaleOutMax  schemas.NumberNoLimits `json:"scale-out-max"`
	Factor       schemas.NumberNoLimits `json:"factor"`
	Offset       schemas.NumberNoLimits `json:"offset"`
}

type AISettings struct {
	DeviceNumber int     `json:"device-number"`
	InputNumber  int     `json:"input-number"`
	IoType       string  `json:"io-type"`
	Decimal      int     `json:"decimal"`
	ScaleEnable  bool    `json:"scale-enable"`
	ScaleInMin   float64 `json:"scale-in-min"`
	ScaleInMax   float64 `json:"scale-in-max"`
	ScaleOutMin  float64 `json:"scale-out-min"`
	ScaleOutMax  float64 `json:"scale-out-max"`
	Factor       float64 `json:"factor"`
	Offset       float64 `json:"offset"`
}

func (inst *AI) buildSchema() *schemas.Schema {
	props := &AISettingsSchema{}

	props.DeviceNumber.Title = "Select IO Device Number (Address)"
	props.DeviceNumber.Default = 1
	props.DeviceNumber.Minimum = 1
	props.DeviceNumber.Maximum = 4

	props.InputNumber.Title = "Select UI Number"
	props.InputNumber.Default = 1
	props.InputNumber.Minimum = 1
	props.InputNumber.Maximum = 8

	props.IoType.Title = "Select UI Input Type"
	props.IoType.Default = string(points.IoTypeVolts)
	props.IoType.Options = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent), string(points.IoTypePulseOnRise), string(points.IoTypePulseOnFall)}
	props.IoType.EnumName = []string{string(points.IoTypeVolts), string(points.IoTypeDigital), string(points.IoTypeTemp), string(points.IoTypeCurrent), string(points.IoTypePulseOnRise), string(points.IoTypePulseOnFall)}

	props.Decimal.Title = "Rounding To # Decimals"
	props.Decimal.Default = 2
	props.Decimal.Minimum = 0
	props.Decimal.Maximum = 10

	props.ScaleEnable.Title = "Enable Scale/Limit Transformation"
	props.ScaleEnable.Default = false

	props.ScaleInMin.Title = "Scale: Input Min"
	props.ScaleInMin.Default = 0

	props.ScaleInMax.Title = "Scale: Input Max"
	props.ScaleInMax.Default = 10

	props.ScaleOutMin.Title = "Scale/Limit: Output Min"
	props.ScaleOutMin.Default = 0

	props.ScaleOutMax.Title = "Scale/Limit: Output Max"
	props.ScaleOutMax.Default = 100

	props.Factor.Title = "Multiplication Factor"
	props.Factor.Default = 0

	props.Offset.Title = "Offset"
	props.Offset.Default = 0

	schema.Set(props)

	uiSchema := array.Map{
		"io-type": array.Map{
			"ui:widget": "select",
		},
		"ui:order": array.Slice{"device-number", "input-number", "io-type", "decimal", "scale-enable", "scale-in-min", "scale-in-max", "scale-out-min", "scale-out-max", "factor", "offset"},
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

func (inst *AI) getSettings(body map[string]interface{}) (*AISettings, error) {
	settings := &AISettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

func (inst *AI) getTransformProps(settings *AISettings) *valueTransformProperties {
	transProps := valueTransformProperties{
		settings.Decimal,
		settings.ScaleEnable,
		settings.ScaleInMin,
		settings.ScaleInMax,
		settings.ScaleOutMin,
		settings.ScaleOutMax,
		settings.Factor,
		settings.Offset,
	}
	return &transProps
}
