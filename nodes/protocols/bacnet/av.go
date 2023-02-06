package bacnetio

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
)

type AV struct {
	*node.Spec
	objectID      points.ObjectID
	objectType    points.ObjectType
	pointUUID     string
	store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

func NewAV(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	body = node.Defaults(body, bacnetAV, category)

	in14 := node.BuildInput(node.In14, node.TypeFloat, nil, body.Inputs, nil)
	in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in14, in15)

	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	currentPriority := node.BuildOutput(node.CurrentPriority, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out, currentPriority)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	flowOptions := &toFlowOptions{}
	node := &AV{
		body,
		0,
		points.AnalogVariable,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
		flowOptions,
	}
	node.SetSchema(node.buildSchema())
	return node, nil
}

func (inst *AV) setObjectId(settings *AVSettings) {
	id := settings.InstanceNumber
	inst.objectID = points.ObjectID(id)
	inst.SetSubTitle(fmt.Sprintf("AV-%d", inst.objectID))
}

func (inst *AV) Process() {
	_, firstLoop := inst.Loop()
	s := inst.GetStore()
	if s == nil {
		return
	}
	if firstLoop {
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		settings, err := inst.getSettings(inst.GetSettings())
		inst.setObjectId(settings)
		transformProps := inst.getTransformProps(settings)
		point := addPoint(points.IoTypeNumber, objectType, inst.objectID, isWriteable, isIO, true, inst.application, transformProps)
		point.Name = inst.GetNodeName()
		point, err = inst.store.AddPoint(point, false)
		if err != nil {
			log.Errorf("bacnet-server add new point type:%s-%d", objectType, inst.objectID)
		}
		s.Set(setUUID(inst.GetParentId(), points.AnalogVariable, inst.objectID), point, 0)
	}

	in14, in15 := fromFlow(inst, inst.objectID)
	pnt := inst.writePointPri(points.AnalogVariable, inst.objectID, in14, in15)
	if pnt != nil {
		inst.WritePinFloat(node.Outp, pnt.PresentValue, 2)
		currentPriority := points.GetHighest(pnt.WriteValue)
		if currentPriority != nil {
			inst.WritePinFloat(node.CurrentPriority, float64(currentPriority.Number), 0)
		}
	} else {
		inst.WritePinNull(node.Outp)
	}
}

func (inst *AV) getPV(objType points.ObjectType, id points.ObjectID) (float64, error) {
	pnt := inst.getPoint(objType, id)
	if pnt != nil {
		return pnt.PresentValue, nil
	}
	return 0, nil
}

func (inst *AV) getPoint(objType points.ObjectType, id points.ObjectID) *points.Point {
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

func (inst *AV) writePointPri(objType points.ObjectType, id points.ObjectID, in14, in15 *float64) *points.Point {
	p := inst.getPoint(objType, id)
	if p == nil {
		return nil
	}

	if p.WriteValueFromBACnet != nil {
		p.WriteValueFromBACnet.P14 = in14
		p.WriteValueFromBACnet.P15 = in15
		p.WriteValue = p.WriteValueFromBACnet
		currentPriority := points.GetHighest(p.WriteValue)
		// fmt.Println(p.WriteValueFromBACnet == nil, "WriteValueFromBACnet", "currentPriority", currentPriority)
		if currentPriority != nil {
			p.PresentValue = currentPriority.Value
		}
		err := inst.updatePoint(objType, id, p)
		if err != nil {
			log.Errorf("error")
		}
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

func (inst *AV) updatePoint(objType points.ObjectType, id points.ObjectID, point *points.Point) error {
	s := inst.GetStore()
	if s == nil {
		return nil
	}
	s.Set(setUUID(inst.GetID(), objType, id), point, 0)
	return nil
}

// Custom Node Settings Schema

type AVSettingsSchema struct {
	InstanceNumber schemas.Integer `json:"instance-number"`
}

type AVSettings struct {
	InstanceNumber int `json:"instance-number"`
}

func (inst *AV) buildSchema() *schemas.Schema {
	props := &AVSettingsSchema{}

	props.InstanceNumber.Title = "Select AV BACnet Instance Number"
	props.InstanceNumber.Default = 1
	props.InstanceNumber.Minimum = 1

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"instance-number"},
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

func (inst *AV) getSettings(body map[string]interface{}) (*AVSettings, error) {
	settings := &AVSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

func (inst *AV) getTransformProps(settings *AVSettings) *ValueTransformProperties {
	transProps := ValueTransformProperties{
		10,
		false,
		0,
		0,
		0,
		0,
		1,
		0,
	}
	return &transProps
}
