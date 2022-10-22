package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
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

//toFlow write the value to the flow, as in an AI write the temp value
func toFlow(body node.Node, objType points.ObjectType, id points.ObjectID, store *points.Store, opts *toFlowOptions) {
	if opts == nil {
		opts = &toFlowOptions{
			precision: 0,
		}
		if objType == points.AnalogInput || objType == points.AnalogOutput {
			opts.precision = 2
		}
	}
	_, v, found := store.GetPresentValueByObject(objType, id) // get the latest value from the point
	if !found {
		log.Error(fmt.Sprintf("bacnet send value to flow runtime failed to find point by object: %s-%d", objType, id))
	}
	body.WritePinFloat(node.Out, v, opts.precision)
}

// mqttPubRunner send messages to the broker, as in read a modbus point and send it to the bacnet server
func (inst *Server) writeRunner() {
	log.Info("start mqtt-pub-runner")
	for {
		for _, point := range inst.store.GetPoints() {
			if inst.store.PendingMQTTPublish(point) {
				err := inst.mqttPublishPV(point)
				if err != nil {
					log.Error(err)
				}
			}
		}
		time.Sleep(runnerDelay * time.Millisecond)
	}
}
