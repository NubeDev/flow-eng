package points

import (
	"github.com/NubeDev/flow-eng/helpers"
)

type Sync struct {
	Points []*point
}

type point struct {
	PointUUID string
	Pending   []*writeSync
}

type writeSync struct {
	UUID        string
	WriteValue  float64
	SyncPending bool
}

func NewSync() *Sync {
	return &Sync{}
}

func (inst *Sync) AddSync(pointUUID string, writeValue float64) {
	p := &point{}
	existing := inst.GetByPoint(pointUUID)
	if existing != nil {
		inst.addWrite(pointUUID, writeValue)
	} else {
		inst.addPoint(pointUUID)
		inst.addWrite(pointUUID, writeValue)
	}
	inst.Points = append(inst.Points, p)
}

func (inst *Sync) addPoint(pointUUID string) {
	p := &point{
		PointUUID: pointUUID,
	}
	inst.Points = append(inst.Points, p)

}

func (inst *Sync) addWrite(pointUUID string, writeValue float64) {
	p := inst.GetByPoint(pointUUID)
	w := &writeSync{
		UUID:        helpers.ShortUUID(),
		WriteValue:  writeValue,
		SyncPending: false,
	}
	p.Pending = append(p.Pending, w)
}

func (inst *Sync) GetPoints() []*point {
	return inst.Points
}

func (inst *Sync) GetByPoint(uuid string) *point {
	for _, p := range inst.Points {
		if p.PointUUID == uuid {
			return p
		}
	}
	return nil
}
