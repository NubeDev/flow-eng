package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *Server) edge28OutputsDispatch() {
	log.Info("start rubix-io-outputs-dispatch")
	for {
		getPoints := inst.store.GetWriteablePointsByApplication(inst.application)
		for _, point := range getPoints { //get the list of the points to update
			if inst.store.PendingWrite(point) {
				if point.ObjectType == points.AnalogOutput {
					_, err := inst.clients.edge28.WriteUO(point)
					if err == nil {
						inst.store.CompletePendingWriteCount(point)
					} else {
						log.Error(err)
					}
				}
				if point.ObjectType == points.BinaryOutput {
					_, err := inst.clients.edge28.WriteDO(point)
					if err == nil {
						inst.store.CompletePendingWriteCount(point)
					} else {
						log.Error(err)
					}
				}
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func (inst *Server) edge28InputsRunner() {
	for {
		analogPoint := inst.store.GetPointsByApplicationAndType(names.Edge, points.AnalogInput)
		analogValues, err := inst.clients.edge28.GetUIs(analogPoint...)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, point := range analogValues {
			inst.store.WriteValueFromRead(point, point.ValueFromRead)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
