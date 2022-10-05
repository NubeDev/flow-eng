package points

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type SyncFrom string // FromMqttPriory, FromFlow, FromRubixIO

const (
	FromRubixIO    SyncFrom = "FromRubixIO"    // write to rubix-io outputs
	FromMqttPriory SyncFrom = "FromMqttPriory" // message from the broker, ie: something wrote via bacnet
	FromFlow       SyncFrom = "FromFlow"       // message from the broker, ie: something wrote via bacnet
)

type SyncList struct {
	Completed bool
}

type writeSync struct {
	UUID        string
	WriteValue  *PriArray
	SyncPending bool
	Time        time.Time
	SyncFrom    SyncFrom
	SyncTo      *SyncList // modbus, rubix-io
}

// CreateSync can come from bacnet or the flow
func (inst *Store) CreateSync(writeValue *PriArray, object ObjectType, id ObjectID, syncFrom SyncFrom, in14, in15 *float64) {
	point := inst.GetPointByObject(object, id)
	if object == "" {
		log.Errorf("bacnet-server: object type type can not be empty")
	}
	if syncFrom == "" {
		log.Errorf("bacnet-server: get sync from can not be empty")
	}
	inst.WritePointValue(point.UUID, writeValue, in14, in15)

}
