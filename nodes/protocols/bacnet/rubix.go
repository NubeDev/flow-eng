package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/rubixio"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *Server) rubixInputsRunner(msg *topics.Message) {
	rubix := &rubixIO.RubixIO{}
	inputsPayload, err := rubix.DecodeInputs(msg.Msg.Payload()) // data comes from mqtt
	if err != nil {
		log.Error(err)
		//return
	}
	for _, point := range inst.store.GetPointsByApplication(names.RubixIO) {
		if point.ObjectType == points.AnalogInput {
			value, err := rubix.DecodeInputValue(point, inputsPayload)
			if err != nil {
				log.Errorf("rubix-io inputs runner: %s", err.Error())
				return
			}
			inst.store.WriteValueFromRead(point, value)
		}
	}
}

func (inst *Server) rubixOutputsDispatch() {
	log.Info("start rubix-io-outputs-dispatch")
	for {
		var pointsToWrite []*points.Point
		getPoints := inst.store.GetWriteablePointsByApplication(inst.application)
		for _, point := range getPoints { //get the list of the points to update
			if inst.store.PendingWrite(point) {
				pointsToWrite = append(pointsToWrite, point)
			}
		}
		if len(pointsToWrite) > 0 {
			bulkPoints, err := inst.clients.rio.BulkWrite(pointsToWrite)
			if err != nil {
				log.Error(err)
			} else {
				for _, point := range bulkPoints {
					inst.store.CompletePendingWriteCount(point)
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
