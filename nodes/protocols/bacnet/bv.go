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

type BV struct {
	*node.Spec
	objectID      points.ObjectID
	objectType    points.ObjectType
	pointUUID     string
	store         *points.Store
	application   names.ApplicationName
	mqttClient    *mqttclient.Client
	toFlowOptions *toFlowOptions
}

func NewBV(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	body = node.Defaults(body, bacnetBV, category)

	in14 := node.BuildInput(node.In14, node.TypeBool, nil, body.Inputs, false)
	in15 := node.BuildInput(node.In15, node.TypeBool, nil, body.Inputs, false)
	inputs := node.BuildInputs(in14, in15)

	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	currentPriority := node.BuildOutput(node.CurrentPriority, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out, currentPriority)

	body = node.BuildNode(body, inputs, outputs, body.Settings)

	flowOptions := &toFlowOptions{}
	n := &BV{
		body,
		0,
		points.BinaryVariable,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
		flowOptions,
	}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *BV) setObjectId(settings *BVSettings) {
	id := settings.InstanceNumber
	inst.objectID = points.ObjectID(id)
	inst.SetSubTitle(fmt.Sprintf("BV-%d", inst.objectID))
}

func (inst *BV) Process() {
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
		s.Set(setUUID(inst.GetParentId(), points.BinaryVariable, inst.objectID), point, 0)
	}

	in14, in15 := fromFlow(inst, inst.objectID)
	pnt := inst.writePointPri(points.BinaryVariable, inst.objectID, in14, in15)
	if pnt != nil {
		inst.WritePinFloat(node.Out, pnt.PresentValue, 2)
		currentPriority := points.GetHighest(pnt.WriteValue)
		if currentPriority != nil {
			inst.WritePinFloat(node.CurrentPriority, float64(currentPriority.Number), 0)
		}
	} else {
		inst.WritePinNull(node.Out)
	}
}

func (inst *BV) getPV(objType points.ObjectType, id points.ObjectID) (float64, error) {
	pnt := inst.getPoint(objType, id)
	if pnt != nil {
		return pnt.PresentValue, nil
	}
	return 0, nil
}

func (inst *BV) getPoint(objType points.ObjectType, id points.ObjectID) *points.Point {
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

func (inst *BV) writePointPri(objType points.ObjectType, id points.ObjectID, in14, in15 *float64) *points.Point {
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

func (inst *BV) updatePoint(objType points.ObjectType, id points.ObjectID, point *points.Point) error {
	s := inst.GetStore()
	if s == nil {
		return nil
	}
	s.Set(setUUID(inst.GetID(), objType, id), point, 0)
	return nil
}

// Custom Node Settings Schema

type BVSettingsSchema struct {
	InstanceNumber schemas.Integer `json:"instance-number"`
}

type BVSettings struct {
	InstanceNumber int `json:"instance-number"`
}

func (inst *BV) buildSchema() *schemas.Schema {
	props := &BVSettingsSchema{}

	props.InstanceNumber.Title = "Select BV Instance Number"
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

func (inst *BV) getSettings(body map[string]interface{}) (*BVSettings, error) {
	settings := &BVSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}

func (inst *BV) getTransformProps(settings *BVSettings) *ValueTransformProperties {
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
