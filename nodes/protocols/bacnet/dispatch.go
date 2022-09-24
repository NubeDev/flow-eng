package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
	"time"
)

/*
DISPATCH
Is where we loop through the store and get the latest write value and then try and write,
the values to protocols like rubix-io, edge28 and modbus

if fail we keep trying but if a new value arrives to the store we will take the latest value,
and disregard the existing
*/

func toFlow(body node.Node, id points.ObjectID) {
	objectType, _, _, err := getBacnetType(body.GetName())
	if err != nil {
		return
	}
	p, v, _ := getStore().GetValueFromReadByObject(objectType, id) // get the latest value from the point store
	body.WritePin(node.Out, v)
	getServer().mqttPublish(p)
}

func (inst *Server) rubixOutputsDispatch() {
	log.Info("start rubix-io-outputs-dispatch")
	for {
		var pointsToWrite []*points.Point
		getPoints := getStore().GetWriteablePointsByApplication(getApplication())
		for _, point := range getPoints { //get the list of the points to update
			sync := getStore().GetLatestSyncValue(point.UUID, points.ToRubixIO)
			if sync != nil {
				point.CurrentSyncUUID = sync.UUID
				pointsToWrite = append(pointsToWrite, point)
			}
		}
		if len(pointsToWrite) > 0 {
			fmt.Println("BULK WIRTE")
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
					//fmt.Println(point.UUID, point.CurrentSyncUUID, updateBacnet)
					getStore().CompleteProtocolWrite(point.UUID, point.CurrentSyncUUID)
					getStore().DeleteSyncWrite(point.UUID, point.CurrentSyncUUID)

					if updateBacnet { // now do it, update bacnet-server

					}
				}
			}
		}
		time.Sleep(runnerDelay * time.Millisecond)
		//time.Sleep(2000 * time.Millisecond)
	}
}
