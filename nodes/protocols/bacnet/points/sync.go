package points

import (
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/names"
	"sort"
	"time"
)

type SyncFrom string // FromMqttPriory, FromFlow, FromRubixIO
type SyncTo string   // FromMqttPriory, FromFlow, FromRubixIO

const (
	ToRubixIOModbus SyncTo = "ToRubixIOModbus" // write to rubix-io outputs
	ToRubixIO       SyncTo = "ToRubixIO"       // write to rubix-io outputs
	ToModbus        SyncTo = "ToModbus"        // write to rubix-io outputs
	ToEdge28        SyncTo = "ToEdge28"        // write to rubix-io outputs
)

const (
	FromRubixIO    SyncFrom = "FromRubixIO"    // write to rubix-io outputs
	FromMqttPriory SyncFrom = "FromMqttPriory" // message from the broker, ie: something wrote via bacnet
	FromFlow       SyncFrom = "FromFlow"       // message from the broker, ie: something wrote via bacnet
)

type SyncList struct {
	Completed bool
	SyncTo    SyncTo
}

type writeSync struct {
	UUID        string
	WriteValue  *PriArray
	SyncPending bool
	Time        time.Time
	SyncFrom    SyncFrom
	SyncTo      []*SyncList // modbus, rubix-io
}

func (inst *Store) AddSync(pointUUID string, writeValue *PriArray, syncFrom SyncFrom, syncTo SyncTo, application names.ApplicationName) {
	p := inst.GetPoint(pointUUID)
	if application == names.RubixIOAndModbus {
		if p != nil {
			modbus := inst.addWrite(writeValue, syncFrom, ToModbus)
			rubix := inst.addWrite(writeValue, syncFrom, ToRubixIO)
			p.Sync = append(p.Sync, modbus, rubix)
		}
	} else {
		if p != nil {
			s := inst.addWrite(writeValue, syncFrom, syncTo)
			p.Sync = append(p.Sync, s)
		}
	}
}

func (inst *Store) GetLatestSyncValue(pointUUID string, to SyncTo) *writeSync {
	s := inst.GetSyncByPoint(pointUUID)
	w := &writeSync{}
	var found bool
	for i, sync := range s {
		for _, list := range sync.SyncTo {
			if i == 0 {
				if list.SyncTo == to {
					w = sync
					found = true
				}
			} else {
				//inst.DeleteSyncWrite(pointUUID, sync.UUID)
			}
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

func (inst *Store) CompleteProtocolWrite(pointUUID, pointCurrentSyncUUID string) bool {
	p := inst.GetSyncByPoint(pointUUID)
	for _, sync := range p {
		if sync.UUID == pointCurrentSyncUUID {
			for _, list := range sync.SyncTo {
				list.Completed = true
				return true
			}
		}
	}
	return false
}

func (inst *Store) addWrite(writeValue *PriArray, syncFrom SyncFrom, syncTo SyncTo) *writeSync {
	to := &SyncList{
		SyncTo: syncTo,
	}
	w := &writeSync{
		UUID:        helpers.ShortUUID(),
		WriteValue:  writeValue,
		SyncPending: false,
		Time:        time.Now(),
		SyncFrom:    syncFrom,
		SyncTo: []*SyncList{
			to,
		},
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
