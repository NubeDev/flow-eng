package bacnetio

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

/*
PROCESS
Is where a message will come from the flow or bacnet
These would only be messages that we need to write to, as in write output-point, to a modbus or edge-28 point
*/

// fromFlow is when a node has been written to from the wire sheet link, as in write a value @16
func fromFlow(body node.Node, objectId points.ObjectID) (*float64, *float64) {
	_, isWriteable, _, err := getBacnetType(body.GetName())
	if err != nil {
		return nil, nil
	}
	var in14 *float64
	var in15 *float64
	if isWriteable {
		in14Value, in14Null := body.ReadPinAsFloat(node.In14)
		if in14Null {
			in14 = nil
		} else {
			in14 = float.New(in14Value)
		}
		in15Value, in15Null := body.ReadPinAsFloat(node.In15)
		if in15Null {
			in15 = nil
		} else {
			in15 = float.New(in15Value)
		}
	}
	if objectId == 0 {
		log.Errorf("bacnet-server: failed to get object-id from node process")
		return nil, nil
	}
	return in14, in15
}
