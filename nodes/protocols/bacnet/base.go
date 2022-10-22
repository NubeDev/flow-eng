package bacnetio

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	points "github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

/*
Scope
example is the rubix-compute-io with 2x IO-16s plugged in

******AI******:
read the UIs from the RC and
address IO-16
- dev-1:UI1 -> AI1
- dev-2:UI1 -> AI9

address RC-IO
- UI1 -> AI17

any UOs and UIs when used as a digital in/out are still an AO or AI
for the edge 28 it has DI/DOs, so we will use the BO/BIs

update point value on bacnet & wire sheet
read the UI and update the wire sheet and update point pv on the BACnet-server via the MQTT broker


******AO******:
Address are the same as above the AI
- read a value from the wire sheet and write device and also the bacnet-server via MQTT

*/

const (
	category   = "bacnet"
	serverNode = "bacnet-server"
	bacnetBI   = "binary-input"
	bacnetBO   = "binary-output"
	bacnetBV   = "binary-variable"
	bacnetAV   = "analog-variable"
	bacnetAO   = "analog-output"
	bacnetAI   = "analog-input"

	typeAI = "ai"
	typeAO = "ao"
	typeAV = "av"
	typeBV = "bv"
)

func getBacnetType(nodeName string) (obj points.ObjectType, isWriteable, isIO bool, err error) {
	switch nodeName {
	case bacnetAI:
		return points.AnalogInput, false, true, nil
	case bacnetAO:
		return points.AnalogOutput, true, true, nil
	case bacnetAV:
		return points.AnalogVariable, true, false, nil
	case bacnetBI:
		return points.BinaryInput, false, true, nil
	case bacnetBO:
		return points.BinaryOutput, true, true, nil
	case bacnetBV:
		return points.BinaryVariable, true, false, nil

	}
	return "", false, false, errors.New(fmt.Sprintf("bacnet add new point object type not found node: %s", nodeName))
}

func nodeDefault(body *node.Spec, nodeName, category string, application names.ApplicationName) (*node.Spec, error) {
	var err error
	body = node.Defaults(body, nodeName, category)
	_, isWriteable, _, err := getBacnetType(nodeName)
	pointName := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	objectIDInput := node.BuildInput(node.ObjectId, node.TypeFloat, 0, body.Inputs)
	var inputs []*node.Input
	if isWriteable {
		if body.GetName() == bacnetBV || body.GetName() == bacnetBO {
			in14 := node.BuildInput(node.In14, node.TypeBool, nil, body.Inputs)
			in15 := node.BuildInput(node.In15, node.TypeBool, nil, body.Inputs)
			inputs = node.BuildInputs(pointName, objectIDInput, in14, in15)
		} else {
			in14 := node.BuildInput(node.In14, node.TypeFloat, nil, body.Inputs)
			in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs)
			inputs = node.BuildInputs(pointName, objectIDInput, in14, in15)
		}

	} else {
		overrideInput := node.BuildInput(node.OverrideInput, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, overrideInput)
	}
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return body, err
}

func addPoint(ioType points.IoType, objectType points.ObjectType, id points.ObjectID, isWriteable, isIO, enable bool) *points.Point {
	point := &points.Point{
		ObjectType:  objectType,
		ObjectID:    id,
		IoType:      ioType,
		IsIO:        isIO,
		IsWriteable: isWriteable,
		Enable:      enable,
	}
	return point

}

// topicBuilder bacnet/ObjectType
func topicObjectBuilder(objectType points.ObjectType) string {
	return fmt.Sprintf("bacnet/%s", objectType)
}

// topicBuilder bacnet/ao/1
func topicBuilder(objectType points.ObjectType, address points.ObjectID) string {
	obj, err := points.ObjectSwitcher(objectType)
	if err != nil {
		log.Error(err)
	}
	return fmt.Sprintf("bacnet/%s/%d", obj, address)
}
