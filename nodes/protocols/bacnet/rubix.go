package bacnet

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/NubeDev/flow-eng/services/rubixio"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *Server) rubixIOBus() {
	if !priorityBus {
		handlerMQTT := bus.Handler{
			Handle: func(ctx context.Context, e bus.Event) {
				go func() {
					decoded := decode(e.Data)
					if decoded != nil {
						inst.rubixInputsRunner(decoded)
					}
				}()
			},
			Matcher: eventbus.RubixIOInputs,
		}
		key := fmt.Sprintf("key_%s", helpers.UUID())
		eventbus.GetBus().RegisterHandler(key, handlerMQTT)
	}
	priorityBus = true
}

func (inst *Server) rubixOutputsRunner() {
	log.Info("start rubix-io-outputs-runner")
	for {
		var pointsToWrite []*points.Point
		for _, point := range getStore().GetWriteablePointsByApplication(applications.RubixIO) { //get the list of the points to update
			sync := getStore().GetLatestSyncValue(point.UUID, points.ToRubixIO)
			if sync != nil {
				pprint.PrintJOSN(sync)
				point.CurrentSyncUUID = sync.UUID
				pointsToWrite = append(pointsToWrite, point)
			}
		}
		if len(pointsToWrite) > 0 {
			bulkPoints, err := inst.rio.BulkWrite(pointsToWrite)
			if err != nil {
				log.Error(err)
			} else {
				for _, point := range bulkPoints {
					syncs := getStore().GetSyncByPoint(point.UUID)
					var updateBacnet bool
					for _, sync := range syncs {
						if sync.SyncFrom == points.FromFlow { // we need to update bacnet server if any of the sync where from the flow
							updateBacnet = true
						}
					}
					getStore().DeleteSyncWrite(point.UUID, point.CurrentSyncUUID)
					if updateBacnet { // now do it, update bacnet-server

					}
				}
			}
		}
		time.Sleep(runnerDelay * time.Millisecond)
	}
}

func (inst *Server) rubixInputsRunner(msg *eventbus.Message) {
	r := &rubixIO.RubixIO{}
	inputs, err := r.DecodeInputs(msg.Msg.Payload())
	if err != nil {
		log.Error(err)
		//return
	}

	for _, point := range getStore().GetPointsByApplication(applications.RubixIO) {
		if point.ObjectType == points.AnalogInput {
			value, err := r.GetInputValue(point, inputs)
			if err != nil {
				return
			}
			fmt.Println(value, "value")
		}
	}

	//fmt.Println("RUBIX-INPUTS", inputs)

}
