package bacnet

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

/*
DISPATCH (POINT-WRITE)
Is where we loop through the store and get the latest write value and then try and write,
the values to protocols like rubix-io, edge28 and modbus

if fail we keep trying but if a new value arrives to the store we will take the latest value,
and disregard the existing
*/

//toFlow write the value to the flow, as in a AI write the temp value
func toFlow(body node.Node, id points.ObjectID, store *points.Store) {
	objectType, _, _, err := getBacnetType(body.GetName())
	if err != nil {
		return
	}
	_, v, _ := store.GetValueFromReadByObject(objectType, id) // get the latest value from the point store
	body.WritePin(node.Out, v)
	//getServer().mqttPublish(p) // MQTT UPDATE
}
