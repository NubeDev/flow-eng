package bacnet

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

const (
	category = "bacnet"
	server   = "server"
	bacnetBI = "binary-input"
	bacnetBO = "binary-output"
	bacnetBV = "binary-variable"
	bacnetAV = "analog-variable"
	bacnetAO = "analog-output"
	bacnetAI = "analog-input"

	typeBV = "bv"
	typeAV = "av"
)

func getBacnetType(nodeName string) (obj bstore.ObjectType, isWriteable, isIO bool, err error) {
	switch nodeName {
	case bacnetAI:
		return bstore.AnalogInput, false, true, nil
	case bacnetAO:
		return bstore.AnalogOutput, true, true, nil
	case bacnetAV:
		return bstore.AnalogVariable, true, false, nil
	case bacnetBI:
		return bstore.BinaryInput, false, true, nil
	case bacnetBO:
		return bstore.BinaryOutput, true, true, nil
	case bacnetBV:
		return bstore.BinaryVariable, true, false, nil

	}
	return "", false, false, errors.New(fmt.Sprintf("bacnet add new point object type not found node: %s", nodeName))
}

func nodeDefault(body *node.Spec, nodeName, category string, application node.ApplicationName) (*node.Spec, error, *bstore.Point) {
	var err error
	body = node.Defaults(body, nodeName, category)

	objectType, isWriteable, isIO, err := getBacnetType(nodeName)

	pointName := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	objectIDInput := node.BuildInput(node.ObjectId, node.TypeFloat, 1, body.Inputs)
	ioType := bstore.IoTypeDigital // TODO make a setting
	enable := true                 // TODO make a setting
	var inputs []*node.Input

	if isWriteable {
		overrideInput := node.BuildInput(node.In16, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, overrideInput)
	} else {
		overrideInput := node.BuildInput(node.OverrideInput, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, overrideInput)
	}

	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))

	parameters := &node.Parameters{
		Application: &node.Application{
			Application: applications.BACnet,
			IsChild:     true,
		},
	}
	body.Parameters = node.BuildParameters(parameters)
	body = node.BuildNode(body, inputs, outputs, nil)

	objectID, _ := objectIDInput.GetValue().(float64)

	point := addPoint(application, ioType, objectType, bstore.ObjectID(objectID), isWriteable, isIO, enable)
	store := getStore()
	point, err = store.AddPoint(point)
	return body, err, point
}

func addPoint(application node.ApplicationName, ioType bstore.IoType, objectType bstore.ObjectType, id bstore.ObjectID, isWriteable, isIO, enable bool) *bstore.Point {
	point := &bstore.Point{
		Application: application,
		ObjectType:  objectType,
		ObjectID:    id,
		IoType:      ioType,
		IsIO:        isIO,
		IsWriteable: isWriteable,
		Enable:      enable,
	}
	return point

}

//func process(body node.Node) {
//	nodeName := body.GetName()
//
//	objectId := body.ReadPin(node.ObjectId)
//	overrideInput := body.ReadPin(node.OverrideInput)
//
//	body.WritePin(node.Out, output)
//}

// topicBuilder bacnet/ao/1
func topicBuilder(objectType string, address bstore.ObjectID) string {
	return fmt.Sprintf("bacnet/%s/%d", objectType, address)
}

// TopicPresentValue bacnet/ao/1/pv
func TopicPresentValue(objectType string, address bstore.ObjectID) string {
	return fmt.Sprintf("%s/pv", topicBuilder(objectType, address))
}

// TopicPriority bacnet/ao/1/pri
func TopicPriority(objectType string, address bstore.ObjectID) string {
	return fmt.Sprintf("%s/pri", topicBuilder(objectType, address))
}
