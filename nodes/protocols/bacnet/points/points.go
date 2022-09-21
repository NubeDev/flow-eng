package points

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
	"math"
)

func (inst *Store) GetPoints() []*Point {
	return inst.Points
}

func (inst *Store) GetPointsByApplication(name node.ApplicationName) []*Point {
	var out []*Point
	for _, point := range inst.GetPoints() {
		if point.Application == name {
			out = append(out, point)
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

func (inst *Store) ReadPresentValue(uuid string) (float64, bool) {
	p := inst.GetPoint(uuid)
	if p != nil {
		return p.ToBacnet, true
	}
	return 0, false
}

func (inst *Store) UpdateBacnetSync(uuid string, value bool) {
	p := inst.GetPoint(uuid)
	p.ToBacnetSyncPending = value
}

func (inst *Store) BacnetSyncPending(uuid string) bool {
	p := inst.GetPoint(uuid)
	return p.ToBacnetSyncPending
}

////GetPointArray get the current priority array
//func (inst *Store) GetPointArray(uuid string) *PriArray {
//	p := inst.GetPoint(uuid)
//	return p.ToBacnet
//}

func cov(existing, new, cov float64) bool {
	v := math.Abs(existing-new) <= cov
	return !v
}

//WritePointValue to is to be written to flow modbus or the wire-sheet @ priority 15
func (inst *Store) WritePointValue(uuid string, value float64) bool {
	p := inst.GetPoint(uuid)
	if p != nil {
		c := cov(p.ToBacnet, value, 0.5)
		if c {
			log.Infof("store write point value type:%s-%d value:%f  uuid:%s", p.ObjectType, p.ObjectID, value, uuid)
			p.ToBacnetSyncPending = c
			p.ToBacnet = value
			return true
		}
	}
	return false
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
