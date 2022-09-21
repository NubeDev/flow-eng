package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

func process(body node.Node) {
	objectType, isWriteable, isIO, err := getBacnetType(body.GetName())
	fmt.Println(objectType, isWriteable, isIO, err)
	if err != nil {
		return
	}
	if isWriteable {
		in1 := body.ReadPin(node.In16)
		val, ok := in1.(float64)
		if !ok {
			return
		}
		point := getStore().GetPointByObject(objectType, 1)
		if point != nil {
			getStore().AddSync(point.UUID, val, points.FromRubixIO, setToSync())
			getStore().WritePointValue(point.UUID, val)

		}

	}

}
