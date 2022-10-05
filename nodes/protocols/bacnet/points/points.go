package points

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"math"
	"reflect"
)

type Point struct {
	UUID               string                `json:"uuid"`
	Application        names.ApplicationName `json:"application"`
	ObjectType         ObjectType            `json:"objectType"`
	ObjectID           ObjectID
	presentValue       *float64
	priAndValue        *priAndValue
	writeValue         float64
	IoType             IoType
	IsIO               bool // if it's an io-pin for a real device
	IsWriteable        bool
	Enable             bool
	ValueFromRead      float64
	WriteValue         *PriArray
	PendingWriteCount  uint64
	PendingMQTTPublish bool
}

func (inst *Store) GetPoints() []*Point {
	return inst.Points
}

func (inst *Store) GetWriteablePointsByApplication(name names.ApplicationName) []*Point {
	var out []*Point
	app := inst.GetApplication()
	var rubix bool
	if app == names.RubixIO || app == names.RubixIOAndModbus {
		rubix = true
	}
	for _, point := range inst.GetPoints() {
		if rubix {
			if point.Application == name {
				if point.IsWriteable {
					out = append(out, point)
				}
			}
		} else {
			if point.Application == name {
				if point.IsWriteable {
					out = append(out, point)
				}
			}
		}
	}
	return out
}

func (inst *Store) GetPointsByApplication(name names.ApplicationName) []*Point {
	var out []*Point
	for _, point := range inst.GetPoints() {
		if point.Application == name {
			out = append(out, point)
		}
	}
	return out
}

func (inst *Store) GetPointsByApplicationAndType(name names.ApplicationName, t ObjectType) []*Point {
	var out []*Point
	for _, point := range inst.GetPoints() {
		if point.Application == name {
			if point.ObjectType == t {
				out = append(out, point)
			}
		}
	}
	return out
}

func (inst *Store) GetPointByObject(t ObjectType, id ObjectID) *Point {
	for _, point := range inst.GetPoints() {
		if point.ObjectType == t {
			if point.ObjectID == id {
				return point
			}
		}
	}
	return nil
}

func (inst *Store) GetPoint(uuid string) *Point {
	for _, point := range inst.GetPoints() {
		if point.UUID == uuid {
			return point
		}
	}
	return nil
}

func cov(existing, new, cov float64) bool {
	v := math.Abs(existing-new) <= cov
	return !v
}

func (inst *Store) mergePriority(p2 *PriArray, in14, in15 *float64) *PriArray {
	if p2 == nil {
		p2 = &PriArray{
			P14: in14, // these are reversed for the flow
			P15: in15, // these are reversed for the flow
		}
		return p2
	}
	out := &PriArray{
		P1:  p2.P1,
		P2:  p2.P2,
		P3:  p2.P3,
		P4:  p2.P4,
		P5:  p2.P5,
		P6:  p2.P6,
		P7:  p2.P7,
		P8:  p2.P8,
		P9:  p2.P9,
		P10: p2.P10,
		P11: p2.P11,
		P12: p2.P12,
		P13: p2.P13,
		P14: in14, // these are reversed for the flow
		P15: in15, // these are reversed for the flow
		P16: p2.P16,
	}
	return out
}

//GetValueFromReadByObject get that value that has already been stored
func (inst *Store) GetValueFromReadByObject(t ObjectType, id ObjectID) (*Point, float64, bool) {
	p := inst.GetPointByObject(t, id)
	if p != nil {
		return p, p.ValueFromRead, true
	}
	return nil, 0, false
}

//GetValueFromRead get that value that has already been stored
func (inst *Store) GetValueFromRead(uuid string) (float64, bool) {
	p := inst.GetPoint(uuid)
	if p != nil {
		return p.ValueFromRead, true
	}
	return 0, false
}

//WriteValueFromRead this is a value from a modbus input or rubix-io input
func (inst *Store) WriteValueFromRead(uuid string, value float64) bool {
	p := inst.GetPoint(uuid)
	if p != nil {
		p.ValueFromRead = value
		return true
	}
	return false
}

//WritePointValue to is to be written to flow modbus or the wire-sheet @ priority 15
func (inst *Store) WritePointValue(uuid string, value *PriArray, in14, in15 *float64) {
	var cov bool
	p := inst.GetPoint(uuid)
	if p != nil {
		if value == nil {
			c := inst.mergePriority(p.WriteValue, in14, in15)
			cov = !reflect.DeepEqual(c, p.WriteValue)
			if cov {
				inst.AddPendingWriteCount(p)
			}
			p.WriteValue = c
		} else {
			c := inst.mergePriority(value, in14, in15)
			p.WriteValue = c
			inst.AddPendingWriteCount(p)
		}
	}
}

func (inst *Store) SetPendingMQTTPublish(point *Point) {
	point.PendingMQTTPublish = true
}

func (inst *Store) PendingMQTTPublish(point *Point) bool {
	return point.PendingMQTTPublish
}

func (inst *Store) CompleteMQTTPublish(point *Point) {
	point.PendingMQTTPublish = false
}

func (inst *Store) PendingWrite(point *Point) bool {
	if point.PendingWriteCount > 0 {
		return true
	}
	return false
}

func (inst *Store) AddPendingWriteCount(point *Point) {
	point.PendingWriteCount++
}
func (inst *Store) CompletePendingWriteCount(point *Point) {
	point.PendingWriteCount--
	inst.SetPendingMQTTPublish(point)
}

func (inst *Store) GetByType(objectType ObjectType) (out []*Point, count int) {
	out = []*Point{}
	for _, pnt := range inst.GetPoints() {
		if pnt.ObjectType == objectType {
			out = append(out, pnt)
		}
	}
	return out, len(out)
}

func (inst *Store) CheckExistingPointErr(point *Point) error {
	if inst.CheckExistingPoint(point) {
		return errors.New(fmt.Sprintf("store-add-point: point is existing object-type:%s:%d", point.ObjectType, point.ObjectID))
	}
	return nil
}

func (inst *Store) CheckExistingPoint(point *Point) bool {
	for _, pnt := range inst.GetPoints() {
		if pnt.ObjectType == point.ObjectType {
			if pnt.ObjectID == point.ObjectID {
				return true
			}
		}
	}
	return false
}
