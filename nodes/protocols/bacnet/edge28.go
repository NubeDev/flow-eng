package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

func (inst *Server) edge28OutputsDispatch() {
	log.Info("start rubix-io-outputs-dispatch")
	for {
		getPoints := inst.store.GetWriteablePointsByApplication(inst.application)
		for _, point := range getPoints { //get the list of the points to update
			sync := inst.store.GetLatestSyncValue(point.UUID, points.ToEdge28)
			if sync != nil {
				point.CurrentSyncUUID = sync.UUID
				if point.ObjectType == points.AnalogOutput {
					inst.clients.edge28.WriteUO(point)
				}
				if point.ObjectType == points.BinaryOutput {
					inst.clients.edge28.WriteDO(point)
				}
				inst.store.CompleteProtocolWrite(point.UUID, point.CurrentSyncUUID)
				inst.store.DeleteSyncWrite(point.UUID, point.CurrentSyncUUID)
			}
		}
	}

}

func (inst *Server) edge28InputsRunner() {
	analogPoint := inst.store.GetPointsByApplicationAndType(names.RubixIO, points.AnalogInput)
	analogValues, err := inst.clients.edge28.GetUIs(analogPoint...)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, point := range analogValues {
		fmt.Println(point.UUID, point.ValueFromRead)
		inst.store.WriteValueFromRead(point.UUID, point.ValueFromRead)
	}
}
