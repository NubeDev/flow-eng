package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/rubixio"
	log "github.com/sirupsen/logrus"
	"time"
)

//var rubixIOBus mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
//	mes := &topics.Message{UUID: helpers.ShortUUID("bus"), Msg: msg}
//	if topics.CheckRubixIO(msg.Topic()) {
//		inst.rubixInputsRunner(mes)
//	}
//
//}

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
			inst.store.WriteValueFromRead(point.UUID, value)
			fmt.Println(value, "rubix-io-input-value")
		}
	}

}

func (inst *Server) rubixOutputsDispatch() {
	log.Info("start rubix-io-outputs-dispatch")
	for {
		var pointsToWrite []*points.Point
		getPoints := inst.store.GetWriteablePointsByApplication(inst.application)
		for _, point := range getPoints { //get the list of the points to update
			sync := inst.store.GetLatestSyncValue(point.UUID, points.ToRubixIO)
			if sync != nil {
				point.CurrentSyncUUID = sync.UUID
				pointsToWrite = append(pointsToWrite, point)
			}
		}
		if len(pointsToWrite) > 0 {
			bulkPoints, err := inst.clients.rio.BulkWrite(pointsToWrite)
			if err != nil {
				log.Error(err)
			} else {
				for _, point := range bulkPoints {
					syncs := inst.store.GetSyncByPoint(point.UUID)
					var updateBacnet bool
					for _, sync := range syncs {
						if sync.SyncFrom == points.FromFlow { // we need to update bacnet server if any of the sync where from the flow
							updateBacnet = true
						}
					}
					//fmt.Println(point.UUID, point.CurrentSyncUUID, updateBacnet)
					inst.store.CompleteProtocolWrite(point.UUID, point.CurrentSyncUUID)
					inst.store.DeleteSyncWrite(point.UUID, point.CurrentSyncUUID)

					if updateBacnet { // now do it, update bacnet-server

					}
				}
			}
		}
		time.Sleep(runnerDelay * time.Millisecond)
		//time.Sleep(2000 * time.Millisecond)
	}
}
