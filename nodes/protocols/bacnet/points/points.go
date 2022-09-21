package points

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
	"math"
)

type Point struct {
	UUID               string               `json:"uuid"`
	Application        node.ApplicationName `json:"application"`
	ObjectType         ObjectType           `json:"objectType"`
	ObjectID           ObjectID             `json:"ObjectID"`
	presentValue       *float64
	priAndValue        *priAndValue
	writeValue         float64
	priArray           *PriArray
	IoType             IoType `json:"ioType"` // temp
	IsIO               bool   // if it's an io-pin for a real device
	IsWriteable        bool
	Enable             bool
	IOWriteSyncPending bool
	WriteSyncPending   bool
	WriteValue         float64
	WriteCOV           float64
	Sync               []*writeSync
	CurrentSyncUUID    string
}

func (inst *Store) GetPoints() []*Point {
	return inst.Points
}

func (inst *Store) GetWriteablePointsByApplication(name node.ApplicationName) []*Point {
	var out []*Point
	for _, point := range inst.GetPoints() {
		if point.Application == name {
			if point.IsWriteable {
				out = append(out, point)
			}

		}
	}
	return out
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
		return p.WriteValue, true
	}
	return 0, false
}

func (inst *Store) UpdateBacnetSync(uuid string, value bool) {
	p := inst.GetPoint(uuid)
	p.WriteSyncPending = value
}

func (inst *Store) BacnetSyncPending(uuid string) bool {
	p := inst.GetPoint(uuid)
	return p.WriteSyncPending
}

func (inst *Store) UpdateIOSync(uuid string, value bool) {
	p := inst.GetPoint(uuid)
	p.IOWriteSyncPending = value
}

func cov(existing, new, cov float64) bool {
	v := math.Abs(existing-new) <= cov
	return !v
}

//WritePointValue to is to be written to flow modbus or the wire-sheet @ priority 15
func (inst *Store) WritePointValue(uuid string, value float64) bool {
	p := inst.GetPoint(uuid)
	if p != nil {
		c := cov(p.WriteValue, value, 0.5)
		if c {
			log.Infof("store write point value type:%s-%d value:%f  uuid:%s", p.ObjectType, p.ObjectID, value, uuid)
			p.WriteSyncPending = true
			p.IOWriteSyncPending = true
			p.WriteValue = value
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
