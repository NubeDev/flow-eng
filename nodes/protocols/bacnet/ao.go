package bacnetio

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
	"strings"
)

type AO struct {
	*node.Spec
	objectID      points.ObjectID
	objectType    points.ObjectType
	pointUUID     string
	store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

func NewAO(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	body = node.Defaults(body, bacnetAO, category)

	in14 := node.BuildInput(node.In14, node.TypeFloat, nil, body.Inputs, false)
	in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs, false)
	inputs := node.BuildInputs(in14, in15)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	currentPriority := node.BuildOutput(node.CurrentPriority, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out, currentPriority)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	flowOptions := &toFlowOptions{}
	n := &AO{
		body,
		0,
		points.AnalogOutput,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
		flowOptions,
	}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *AO) setObjectId(settings *AOSettings) {
	devNum := settings.DeviceNumber
	uoNum := settings.InputNumber
	id := ((devNum - 1) * 8) + uoNum
	inst.objectID = points.ObjectID(id)
	name := bacnetAddress(4, "AO", "UO")
	if len(name) >= id {
		if settings != nil {
			ioType := strings.ReplaceAll(settings.IoType, "_", " ")
			inst.SetSubTitle(strings.ToUpper(fmt.Sprintf("%s %s", name[id-1], ioType)))
		} else {
			inst.SetSubTitle(name[id-1])
		}
	}
}

func (inst *AO) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())
	// transformProps := inst.getTransformProps(settings)
	transformProps := &ValueTransformProperties{}
	_, firstLoop := inst.Loop()
	s := inst.GetStore()
	if s == nil {
		return
	}
	if firstLoop {
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		inst.setObjectId(settings)
		ioType := settings.IoType
		if ioType == "" {
			ioType = string(points.IoTypeVolts)
		}
		point := addPoint(points.IoType(ioType), objectType, inst.objectID, isWriteable, isIO, true, inst.application, transformProps)
		point.Name = inst.GetNodeName()
		point, err = inst.store.AddPoint(point, false)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
		s.Set(setUUID(inst.GetParentId(), points.AnalogOutput, inst.objectID), point, 0)
	}

	in14, in15 := fromFlow(inst, inst.objectID)
	pnt := inst.writePointPri(points.AnalogOutput, inst.objectID, in14, in15)
	if pnt != nil {
		value := modbusScaleOutput(pnt.PresentValue, pnt.Offset)
		inst.WritePinFloat(node.Out, value, 2)
		currentPriority := points.GetHighest(pnt.WriteValue)
		if currentPriority != nil {
			inst.WritePinFloat(node.CurrentPriority, float64(currentPriority.Number), 0)
		}
	} else {
		inst.WritePinNull(node.Out)
	}
}

func scaleAO(value float64, isBO bool) float64 {
	if isBO {
		if value > 0 {
			return 10
		} else {
			return 0
		}
	}
	return float.Scale(value, 0, 100, 0, 10)
}

func (inst *AO) getPV(objType points.ObjectType, id points.ObjectID) (float64, error) {
	pnt := inst.getPoint(objType, id)
	if pnt != nil {
		return pnt.PresentValue, nil
	}
	return 0, nil
}

func (inst *AO) getPoint(objType points.ObjectType, id points.ObjectID) *points.Point {
	s := inst.GetStore()
	if s == nil {
		return nil
	}
	d, ok := s.Get(setUUID(inst.ParentId, objType, id))
	if ok {
		return d.(*points.Point)
	}
	return nil
}

func (inst *AO) writePointPri(objType points.ObjectType, id points.ObjectID, in14, in15 *float64) *points.Point {
	p := inst.getPoint(objType, id)
	if p == nil {
		return nil
	}
	if p.WriteValueFromBACnet != nil {
		p.WriteValueFromBACnet.P14 = in14
		p.WriteValueFromBACnet.P15 = in15
		p.WriteValue = p.WriteValueFromBACnet
		currentPriority := points.GetHighest(p.WriteValue)
		if currentPriority != nil {
			p.PresentValue = currentPriority.Value
		}
		inst.updatePoint(objType, id, p)
		return p
	} else {
		if p.WriteValue == nil {
			p.WriteValue = points.NewPriArray(in14, in15)
		} else {
			p.WriteValue.P14 = in14
			p.WriteValue.P15 = in15
		}
		currentPriority := points.GetHighest(p.WriteValue)
		if currentPriority != nil {
			p.PresentValue = currentPriority.Value
		}
		inst.updatePoint(objType, id, p)
		return p
	}

}

func (inst *AO) updatePoint(objType points.ObjectType, id points.ObjectID, point *points.Point) error {
	s := inst.GetStore()
	if s == nil {
		return nil
	}
	s.Set(setUUID(inst.GetID(), objType, id), point, 0)
	return nil
}

// Custom Node Settings Schema

type AOSettingsSchema struct {
	DeviceNumber schemas.Integer    `json:"device-number"`
	InputNumber  schemas.Integer    `json:"input-number"`
	IoType       schemas.EnumString `json:"io-type"`
	/*
		Decimal      schemas.Number         `json:"decimal"`
		ScaleEnable  schemas.Boolean        `json:"scale-enable"`
		ScaleInMin   schemas.NumberNoLimits `json:"scale-in-min"`
		ScaleInMax   schemas.NumberNoLimits `json:"scale-in-max"`
		ScaleOutMin  schemas.NumberNoLimits `json:"scale-out-min"`
		ScaleOutMax  schemas.NumberNoLimits `json:"scale-out-max"`
		Factor       schemas.NumberNoLimits `json:"factor"`
		Offset       schemas.NumberNoLimits `json:"offset"`
	*/
}

type AOSettings struct {
	DeviceNumber int    `json:"device-number"`
	InputNumber  int    `json:"input-number"`
	IoType       string `json:"io-type"`
	/*
		Decimal      int     `json:"decimal"`
		ScaleEnable  bool    `json:"scale-enable"`
		ScaleInMin   float64 `json:"scale-in-min"`
		ScaleInMax   float64 `json:"scale-in-max"`
		ScaleOutMin  float64 `json:"scale-out-min"`
		ScaleOutMax  float64 `json:"scale-out-max"`
		Factor       float64 `json:"factor"`
		Offset       float64 `json:"offset"`
	*/
}

func (inst *AO) buildSchema() *schemas.Schema {
	props := &AOSettingsSchema{}

	props.DeviceNumber.Title = "Select IO Device Number (Address)"
	props.DeviceNumber.Default = 1
	props.DeviceNumber.Minimum = 1
	props.DeviceNumber.Maximum = 4

	props.InputNumber.Title = "Select UO Number"
	props.InputNumber.Default = 1
	props.InputNumber.Minimum = 1
	props.InputNumber.Maximum = 8

	props.IoType.Title = "Select UI Input Type"
	props.IoType.Default = string(points.IoTypeVolts)
	props.IoType.Options = []string{string(points.IoTypeVolts), string(points.IoTypeDigital)}
	props.IoType.EnumName = []string{string(points.IoTypeVolts), string(points.IoTypeDigital)}

	/*
		props.Decimal.Title = "Rounding To # Decimals"
		props.Decimal.Default = 2
		props.Decimal.Minimum = 0
		props.Decimal.Maximum = 10

		props.ScaleEnable.Title = "Enable Scale/Limit Transformation"
		props.ScaleEnable.Default = false

		props.ScaleInMin.Title = "Scale: Input Min (Input 0-10v)"
		props.ScaleInMin.Default = 0
		props.ScaleInMin.ReadOnly = true

		props.ScaleInMax.Title = "Scale: Input Max (Input 0-10v)"
		props.ScaleInMax.Default = 10
		props.ScaleInMax.ReadOnly = true

		props.ScaleOutMin.Title = "Scale/Limit: Output Min"
		props.ScaleOutMin.Default = 0

		props.ScaleOutMax.Title = "Scale/Limit: Output Max"
		props.ScaleOutMax.Default = 10

		props.Factor.Title = "Multiplication Factor"
		props.Factor.Default = 0

		props.Offset.Title = "Offset"
		props.Offset.Default = 0

	*/

	schema.Set(props)

	uiSchema := array.Map{
		"io-type": array.Map{
			"ui:widget": "select",
		},
		// "ui:order": array.Slice{"device-number", "input-number", "io-type", "decimal", "scale-enable", "scale-in-min", "scale-in-max", "scale-out-min", "scale-out-max", "factor", "offset"},
		"ui:order": array.Slice{"device-number", "input-number", "io-type"},
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

func (inst *AO) getSettings(body map[string]interface{}) (*AOSettings, error) {
	settings := &AOSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

/*
func (inst *AO) getTransformProps(settings *AOSettings) *ValueTransformProperties {
	transProps := ValueTransformProperties{
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

*/
