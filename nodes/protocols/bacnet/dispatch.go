package bacnetio

import (
	log "github.com/sirupsen/logrus"
	"time"
)

/*
DISPATCH (POINT-WRITE)
Is where we loop through the store and get the latest write value and then try and write,
the values to protocols like rubix-io, edge28 and modbus

if fail we keep trying but if a new value arrives to the store we will take the latest value,
and disregard the existing
*/

type toFlowOptions struct {
	precision int
}

// mqttPubRunner send messages to the broker, as in read a modbus point and send it to the bacnet server
func (inst *Server) writeRunner() {
	log.Info("start mqtt-pub-runner")
	for {
		p, _ := inst.getPoints()
		for _, point := range p {
			//pprint.PrintJOSN(point)
			inst.mqttPublishPV(point)
		}
		time.Sleep(runnerDelay * time.Millisecond)
	}
}
