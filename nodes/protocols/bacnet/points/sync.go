package points

import (
	"github.com/NubeDev/flow-eng/helpers"
	"sort"
	"time"
)

type SyncFrom string // FromMqttPriory, FromFlow, FromRubixIO
type SyncTo string   // FromMqttPriory, FromFlow, FromRubixIO

const (
	ToRubixIO SyncTo = "ToRubixIO" // write to rubix-io outputs
)

const (
	FromRubixIO    SyncFrom = "FromRubixIO"    // write to rubix-io outputs
	FromMqttPriory SyncFrom = "FromMqttPriory" // message from the broker, ie: something wrote via bacnet
	FromFlow       SyncFrom = "FromFlow"       // message from the broker, ie: something wrote via bacnet
)

type writeSync struct {
	UUID        string
	WriteValue  float64
	SyncPending bool
	Time        time.Time
	SyncFrom    SyncFrom
	SyncTo      SyncTo
}

func (inst *Store) AddSync(pointUUID string, writeValue float64, syncFrom SyncFrom, syncTo SyncTo) {
	p := inst.GetPoint(pointUUID)
	if p != nil {
		s := inst.addWrite(writeValue, syncFrom, syncTo)
		p.Sync = append(p.Sync, s)
	}
}

func (inst *Store) GetLatestSyncValue(pointUUID string, to SyncTo) *writeSync {
	s := inst.GetSyncByPoint(pointUUID)
	w := &writeSync{}
	var found bool
	for i, sync := range s {
		if sync.SyncTo == to {
			if i == 0 {
				w = sync
			}
			found = true
			inst.DeleteSyncWrite(pointUUID, sync.UUID)
		}
	}
	if found {
		return w
	}
	return nil
}

func (inst *Store) GetSyncByPoint(pointUUID string) []*writeSync {
	p := inst.GetPoint(pointUUID)
	if p != nil {
		sort.Slice(p.Sync, func(i, j int) bool {
			return p.Sync[i].Time.After(p.Sync[j].Time)
		})
		return p.Sync
	}
	return nil
}

func (inst *Store) addWrite(writeValue float64, syncFrom SyncFrom, syncTo SyncTo) *writeSync {
	w := &writeSync{
		UUID:        helpers.ShortUUID(),
		WriteValue:  writeValue,
		SyncPending: false,
		Time:        time.Now(),
		SyncFrom:    syncFrom,
		SyncTo:      syncTo,
	}
	return w
}

func removeWrite(slice []*writeSync, s int) []*writeSync {
	return append(slice[:s], slice[s+1:]...)
}

func (inst *Store) DeleteSyncWrite(pointUUID, uuid string) bool {
	for _, p := range inst.GetPoints() {
		if p.UUID == pointUUID { // get the point
			for i, sync := range p.Sync {
				if sync.UUID == uuid { // get the sync to delete
					p.Sync = removeWrite(p.Sync, i)
					return true
				}
			}
		}
	}
	return false
}
