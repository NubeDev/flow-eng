package points

type SyncFrom string // FromMqttPriory, FromFlow, FromRubixIO

const (
	FromRubixIO    SyncFrom = "FromRubixIO"    // write to rubix-io outputs
	FromMqttPriory SyncFrom = "FromMqttPriory" // message from the broker, ie: something wrote via bacnet
	FromFlow       SyncFrom = "FromFlow"       // message from the broker, ie: something wrote via bacnet
)

//// CreateSync can come from bacnet or the flow
//func (inst *Store) CreateSync(writeValue *PriArray, object ObjectType, id ObjectID, syncFrom SyncFrom, in14, in15 *float64) {
//	point := inst.GetPointByObject(object, id)
//	if object == "" {
//		log.Errorf("bacnet-server: object type type can not be empty")
//	}
//	if syncFrom == "" {
//		log.Errorf("bacnet-server: get sync from can not be empty")
//	}
//	inst.WritePointValue(point, writeValue, in14, in15, syncFrom)
//
//}
