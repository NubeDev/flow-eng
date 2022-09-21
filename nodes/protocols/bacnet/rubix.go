package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/eventbus"
	"github.com/NubeDev/flow-eng/services/rubixio"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *Server) rubixOutputsRunner() {
	log.Info("start rubix-io-outputs-runner")
	for {
		var pointsToWrite []*points.Point
		for _, point := range getStore().GetWriteablePointsByApplication(applications.RubixIO) { //get the list of the points to update
			if point.IOWriteSyncPending {
				pointsToWrite = append(pointsToWrite, point)
			} else {
				//log.Infof("mqtt-runner-publish point skip as non cov type:%s-%d", point.ObjectType, point.ObjectID)
			}
		}
		if len(pointsToWrite) > 0 {
			err := inst.rio.BulkWrite(pointsToWrite)
			if err != nil {
				log.Error(err)
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
