package bacnetio

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
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
	var err error
	body, err = nodeDefault(body, bacnetAO, category, opts.Application)
	body.SetSchema(buildSchemaUO())
	flowOptions := &toFlowOptions{}
	return &AO{
		body,
		0,
		points.AnalogOutput,
		"",
		opts.Store,
		opts.Application,
		opts.MqttClient,
		flowOptions,
	}, err
}

func (inst *AO) setObjectId() {
	id, _ := inst.ReadPinAsInt(node.ObjectId)
	inst.objectID = points.ObjectID(id)
}
func (inst *AO) Process() {
	_, firstLoop := inst.Loop()
	s := inst.GetStore()
	if s == nil {
		return
	}
	if firstLoop {
		inst.setObjectId()
		objectType, isWriteable, isIO, err := getBacnetType(inst.Info.Name)
		settings, err := getSettings(inst.GetSettings())
		ioType := settings.Io
		if ioType == "" {
			ioType = string(points.IoTypeVolts)
		}
		inst.toFlowOptions.precision = settings.Decimal
		point := addPoint(points.IoType(ioType), objectType, inst.objectID, isWriteable, isIO, true, inst.application)
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
		inst.WritePinFloat(node.Out, pnt.PresentValue, 2)
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
